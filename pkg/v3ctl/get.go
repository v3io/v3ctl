package v3ctl

import (
	"fmt"
	"os"
	"strconv"

	"github.com/v3io/v3io-go/pkg/dataplane"
	"github.com/v3io/v3io-go/pkg/dataplane/http"

	"github.com/nuclio/errors"
	"github.com/nuclio/renderer"
	"github.com/spf13/cobra"
)

type getCommandeer struct {
	cmd            *cobra.Command
	rootCommandeer *RootCommandeer
}

func newGetCommandeer(rootCommandeer *RootCommandeer) *getCommandeer {
	commandeer := &getCommandeer{
		rootCommandeer: rootCommandeer,
	}

	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get resources",
	}

	getContainersCommand := newGetContainersCommandeer(commandeer).cmd
	getStreamsCommand := newGetStreamCommandeer(commandeer).cmd

	cmd.AddCommand(
		getContainersCommand,
		getStreamsCommand,
	)

	commandeer.cmd = cmd

	return commandeer
}

type getContainersCommandeer struct {
	*getCommandeer
}

func newGetContainersCommandeer(getCommandeer *getCommandeer) *getContainersCommandeer {
	commandeer := &getContainersCommandeer{
		getCommandeer: getCommandeer,
	}

	cmd := &cobra.Command{
		Use:   "containers",
		Short: "Get data containers",
		RunE: func(cmd *cobra.Command, args []string) error {

			// initialize root
			if err := getCommandeer.rootCommandeer.initialize(); err != nil {
				return errors.Wrap(err, "Failed to initialize root")
			}

			getContainersInput := v3io.GetContainersInput{}
			getContainersInput.AuthenticationToken = v3iohttp.GenerateAuthenticationToken(getCommandeer.rootCommandeer.username, getCommandeer.rootCommandeer.password)
			getContainersInput.AccessKey = getCommandeer.rootCommandeer.accessKey

			response, err := getCommandeer.rootCommandeer.dataPlaneContext.GetContainersSync(&getContainersInput)
			if err != nil {
				return errors.Wrap(err, "Failed to get containers")
			}

			commandeer.renderContainers(response.Output.(*v3io.GetContainersOutput).Results.Containers)

			return nil
		},
	}

	commandeer.cmd = cmd

	return commandeer
}

func (c *getCommandeer) renderContainers(containerInfos []v3io.ContainerInfo) {
	renderer := renderer.NewRenderer(os.Stdout)

	switch c.rootCommandeer.output {
	case "", "text":

		var records [][]string
		for _, containerInfo := range containerInfos {
			records = append(records, []string{
				strconv.Itoa(containerInfo.ID),
				containerInfo.Name,
				containerInfo.CreationDate,
			})
		}

		renderer.RenderTable([]string{"ID", "Name", "Creation date"}, records)
	case "yaml":
		renderer.RenderYAML(containerInfos) // nolint: errcheck
	case "json":
		renderer.RenderJSON(containerInfos) // nolint: errcheck
	}
}

type getStreamsCommandeer struct {
	*getCommandeer
}

func newGetStreamCommandeer(getCommandeer *getCommandeer) *getStreamsCommandeer {
	commandeer := &getStreamsCommandeer{
		getCommandeer: getCommandeer,
	}

	cmd := &cobra.Command{
		Use:   "stream path",
		Short: "Get all streams at a given path",
		RunE: func(cmd *cobra.Command, args []string) error {

			// if we got positional arguments
			if len(args) != 1 {
				return errors.New("Stream get requires a stream path")
			}

			// initialize root
			if err := getCommandeer.rootCommandeer.initialize(); err != nil {
				return errors.Wrap(err, "Failed to initialize root")
			}

			response, err := getCommandeer.rootCommandeer.dataPlaneContext.GetContainerContentsSync(&v3io.GetContainerContentsInput{
				Path: args[0],
			})

			if err != nil {
				return errors.Wrapf(err, "Failed to get container contents at %s", args[0])
			}

			for _, content := range response.Output.(*v3io.GetContainerContentsOutput).Contents {
				fmt.Println(content.Key)
			}

			return nil
		},
	}

	commandeer.cmd = cmd

	return commandeer
}
