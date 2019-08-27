package v3ctl

import (
	"os"
	"sync"

	"github.com/nuclio/errors"
	"github.com/nuclio/logger"
	"github.com/nuclio/zap"
	"github.com/spf13/cobra"
	"github.com/v3io/v3io-go/pkg/controlplane"
	"github.com/v3io/v3io-go/pkg/controlplane/http"
	"github.com/v3io/v3io-go/pkg/dataplane"
	"github.com/v3io/v3io-go/pkg/dataplane/http"
)

type RootCommandeer struct {
	Logger              logger.Logger
	cmd                 *cobra.Command
	verbose             bool
	server              string
	controlServer       string
	containerName       string
	logLevel            string
	username            string
	password            string
	accessKey           string
	controlPlaneSession v3ioc.Session
	dataPlaneContext    v3io.Context
	container           v3io.Container
	v3ioLock            sync.Mutex
	output              string
}

// Execute uses os.Args to execute the command
func (rc *RootCommandeer) Execute() error {
	return rc.cmd.Execute()
}

func NewRootCommandeer() *RootCommandeer {
	commandeer := &RootCommandeer{}

	cmd := &cobra.Command{
		Use:           "v3ctl [command]",
		Short:         "v3io command-line interface",
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	defaultV3ioServer := os.Getenv("V3IO_API")

	cmd.PersistentFlags().StringVarP(&commandeer.logLevel, "log-level", "v", "info",
		`Verbose output. Add "=<level>" to set the log level -
debug | info | warn | error. For example: -v=warn`)
	cmd.PersistentFlags().StringVarP(&commandeer.server, "server", "s", defaultV3ioServer,
		`Web-gateway (web-APIs) service endpoint of an instance of
the Iguazio Continuous Data Platform, of the format
"<IP address>:<port number=8081>". Examples: "localhost:8081"
(when running on the target platform); "192.168.1.100:8081".`)
	cmd.PersistentFlags().StringVarP(&commandeer.controlServer, "control-server", "", defaultV3ioServer,
		`Service endpoint of the control server`)
	cmd.PersistentFlags().StringVarP(&commandeer.containerName, "container", "c", "",
		`The name of an Iguazio Continuous Data Platform data container
in which to create the TSDB table. Example: "bigdata".`)
	cmd.PersistentFlags().StringVarP(&commandeer.username, "username", "u", "",
		"Username of an Iguazio Continuous Data Platform user.")
	cmd.PersistentFlags().StringVarP(&commandeer.password, "password", "p", "",
		"Password of the configured user (see -u|--username).")
	cmd.PersistentFlags().StringVarP(&commandeer.accessKey, "access-key", "k", "",
		`Access-key for accessing the required table.
If access-key is passed, it will take precedence on user/password authentication.`)
	cmd.PersistentFlags().StringVarP(&commandeer.output, "output", "o", "text",
		`One of text, wide, yaml, json`)

	// add children
	cmd.AddCommand(
		newCreateCommandeer(commandeer).cmd,
		newGetCommandeer(commandeer).cmd,
		newDeleteCommandeer(commandeer).cmd,
		newLsCommandeer(commandeer).cmd,
	)

	commandeer.cmd = cmd

	return commandeer
}

func (rc *RootCommandeer) initialize() error {
	var err error

	rc.Logger, err = rc.createLogger()
	if err != nil {
		return errors.Wrap(err, "Failed to create logger")
	}

	rc.dataPlaneContext, err = v3iohttp.NewContext(rc.Logger, &v3io.NewContextInput{ClusterEndpoints: []string{rc.server}})
	if err != nil {
		return errors.Wrap(err, "Failed to create v3io context")
	}

	if rc.containerName == "" {
		return errors.New("Container must be specified (use --container)")
	}

	var username string
	var password string
	var accessKey = rc.accessKey

	// Only use username and password from cli if access key was not also provided on cli.
	if accessKey == "" {
		username = rc.username
		password = rc.password

		// Only use V3IO_ACCESS_KEY if no credentials were provided on cli.
		if password == "" {
			accessKey = os.Getenv("V3IO_ACCESS_KEY")
		}
	}

	session, err := rc.dataPlaneContext.NewSession(&v3io.NewSessionInput{
		Username:  username,
		Password:  password,
		AccessKey: accessKey,
	})

	if err != nil {
		return errors.Wrap(err, "Failed to create session")
	}

	rc.container, err = session.NewContainer(&v3io.NewContainerInput{
		ContainerName: rc.containerName,
	})

	if err != nil {
		return errors.Wrap(err, "Failed to open container")
	}

	return nil
}

func (rc *RootCommandeer) createLogger() (logger.Logger, error) {
	var loggerLevel nucliozap.Level

	if rc.logLevel == "debug" {
		loggerLevel = nucliozap.DebugLevel
	} else {
		loggerLevel = nucliozap.InfoLevel
	}

	loggerInstance, err := nucliozap.NewNuclioZapCmd("v3ctl", loggerLevel)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create logger")
	}

	return loggerInstance, nil
}

func (rc *RootCommandeer) getControlPlaneSession() (v3ioc.Session, error) {
	var err error

	rc.v3ioLock.Lock()
	defer rc.v3ioLock.Unlock()

	if rc.controlPlaneSession != nil {
		return rc.controlPlaneSession, nil
	}

	createSessionInput := v3ioc.NewSessionInput{}
	createSessionInput.Endpoints = []string{rc.controlServer}
	createSessionInput.Username = rc.username
	createSessionInput.Password = rc.password
	// createSessionInput.AccessKey = rc.accessKey

	rc.controlPlaneSession, err = v3iochttp.NewSession(rc.Logger, &createSessionInput)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create control plane session")
	}

	return rc.controlPlaneSession, nil
}
