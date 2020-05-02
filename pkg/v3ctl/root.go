package v3ctl

import (
	"os"
	"strings"
	"sync"

	"github.com/nuclio/errors"
	"github.com/nuclio/logger"
	"github.com/nuclio/loggerus"
	"github.com/nuclio/renderer"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/v3io/v3io-go/pkg/controlplane"
	v3iochttp "github.com/v3io/v3io-go/pkg/controlplane/http"
	"github.com/v3io/v3io-go/pkg/dataplane"
	"github.com/v3io/v3io-go/pkg/dataplane/http"
)

type RootCommandeer struct {
	Logger              logger.Logger
	WebapiURL           string
	ControlURL          string
	ContainerName       string
	logLevel            string
	Username            string
	Password            string
	AccessKey           string
	ControlPlaneSession v3ioc.Session
	DataPlaneContext    v3io.Context
	Container           v3io.Container
	Output              string

	cmd      *cobra.Command
	v3ioLock sync.Mutex
}

func NewRootCommandeer() (*RootCommandeer, error) {
	commandeer := &RootCommandeer{}

	cmd := &cobra.Command{
		Use:           "v3ctl [command]",
		Short:         "v3io command-line interface",
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	defaultV3ioServer := os.Getenv("V3IO_API")
	if !strings.HasPrefix(defaultV3ioServer, "http") {
		defaultV3ioServer = "http://" + defaultV3ioServer
	}

	cmd.PersistentFlags().StringVarP(&commandeer.logLevel, "log-level", "v", "info",
		`Verbose output. Add "=<level>" to set the log level -
debug | info | warn | error. For example: -v=warn`)
	cmd.PersistentFlags().StringVarP(&commandeer.WebapiURL, "webapi-url", "w", defaultV3ioServer,
		`Web-gateway (web-APIs) service endpoint of an instance of
the Iguazio Continuous Data Platform, of the format
"<IP address>:<port number=8081>". Examples: "localhost:8081"
(when running on the target platform); "192.168.1.100:8081".`)
	cmd.PersistentFlags().StringVarP(&commandeer.ControlURL, "control-url", "", defaultV3ioServer,
		`Service endpoint of the control server`)
	cmd.PersistentFlags().StringVarP(&commandeer.ContainerName, "container", "c", "",
		`The name of an Iguazio Continuous Data Platform data container
in which to create the TSDB table. Example: "bigdata".`)
	cmd.PersistentFlags().StringVarP(&commandeer.Username, "username", "u", "",
		"Username of an Iguazio Continuous Data Platform user.")
	cmd.PersistentFlags().StringVarP(&commandeer.Password, "password", "p", "",
		"Password of the configured user (see -u|--username).")
	cmd.PersistentFlags().StringVarP(&commandeer.AccessKey, "access-key", "k", "",
		`Access-key for accessing the required table.
If access-key is passed, it will take precedence on user/password authentication.`)
	cmd.PersistentFlags().StringVarP(&commandeer.Output, "output", "o", "text",
		`One of text, wide, yaml, json`)

	createCommandeerInstance, err := newCreateCommandeer(commandeer)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create create commandeer")
	}

	deleteCommandeerInstance, err := newDeleteCommandeer(commandeer)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create delete commandeer")
	}

	getCommandeerInstance, err := newGetCommandeer(commandeer)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create get commandeer")
	}

	cmd.AddCommand(createCommandeerInstance.Cmd)
	cmd.AddCommand(deleteCommandeerInstance.Cmd)
	cmd.AddCommand(getCommandeerInstance.Cmd)

	commandeer.cmd = cmd

	return commandeer, nil
}

// Execute uses os.Args to execute the command
func (c *RootCommandeer) Execute() error {
	return c.cmd.Execute()
}

func (c *RootCommandeer) GetControlPlaneSession() (v3ioc.Session, error) {
	var err error

	c.v3ioLock.Lock()
	defer c.v3ioLock.Unlock()

	if c.ControlPlaneSession != nil {
		return c.ControlPlaneSession, nil
	}

	createSessionInput := v3ioc.NewSessionInput{}
	createSessionInput.Endpoints = []string{c.ControlURL}
	createSessionInput.Username = c.Username
	createSessionInput.Password = c.Password
	// createSessionInput.AccessKey = c.accessKey

	c.ControlPlaneSession, err = v3iochttp.NewSession(c.Logger, &createSessionInput)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create control plane session")
	}

	return c.ControlPlaneSession, nil
}

func (c *RootCommandeer) Initialize() error {
	var err error

	c.Logger, err = c.createLogger()
	if err != nil {
		return errors.Wrap(err, "Failed to create logger")
	}

	c.DataPlaneContext, err = v3iohttp.NewContext(c.Logger,
		&v3iohttp.NewContextInput{})
	if err != nil {
		return errors.Wrap(err, "Failed to create v3io context")
	}

	if c.ContainerName == "" {
		return errors.New("Container must be specified (use --container)")
	}

	var username string
	var password string
	var accessKey = c.AccessKey

	// Only use username and password from cli if access key was not also provided on cli.
	if accessKey == "" {
		username = c.Username
		password = c.Password

		// Only use V3IO_ACCESS_KEY if no credentials were provided on cli.
		if password == "" {
			accessKey = os.Getenv("V3IO_ACCESS_KEY")
		}
	}

	session, err := c.DataPlaneContext.NewSession(&v3io.NewSessionInput{
		URL:       c.WebapiURL,
		Username:  username,
		Password:  password,
		AccessKey: accessKey,
	})

	if err != nil {
		return errors.Wrap(err, "Failed to create session")
	}

	c.Container, err = session.NewContainer(&v3io.NewContainerInput{
		ContainerName: c.ContainerName,
	})

	if err != nil {
		return errors.Wrap(err, "Failed to open container")
	}

	return nil
}

func (c *RootCommandeer) Render(info interface{}, columns []string, records [][]string) error {
	renderer := renderer.NewRenderer(os.Stdout)

	switch c.Output {
	case "", "text":
		if len(columns) == 0 {
			return errors.New("Table render not supported. Try json or yaml")
		}

		renderer.RenderTable(columns, records)
	case "yaml":
		return renderer.RenderYAML(info) // nolint: errcheck
	case "json":
		return renderer.RenderJSON(info) // nolint: errcheck
	}

	return nil
}

func (c *RootCommandeer) createLogger() (logger.Logger, error) {
	var loggerLevel logrus.Level

	if c.logLevel == "debug" {
		loggerLevel = logrus.DebugLevel
	} else {
		loggerLevel = logrus.InfoLevel
	}

	return loggerus.NewTextLoggerus("v3ctl", loggerLevel, os.Stdout, false)
}
