package container

import (
	"os"
	"strconv"

	"github.com/v3io/v3ctl/pkg/v3ctl"

	"github.com/nuclio/errors"
	"github.com/nuclio/renderer"
	"github.com/spf13/cobra"
	v3io "github.com/v3io/v3io-go/pkg/dataplane"
	v3iohttp "github.com/v3io/v3io-go/pkg/dataplane/http"
)

type getContainerCommandeer struct {
	*v3ctl.GetCommandeer
}

func newGetContainerCommandeer(getCommandeer *v3ctl.GetCommandeer) (*getContainerCommandeer, error) {
	commandeer := &getContainerCommandeer{
		GetCommandeer: getCommandeer,
	}

	cmd := &cobra.Command{
		Use:   "container name",
		Short: "Get a data container",
		RunE: func(cmd *cobra.Command, args []string) error {

			// initialize root
			if err := getCommandeer.RootCommandeer.Initialize(); err != nil {
				return errors.Wrap(err, "Failed to initialize root")
			}

			getContainersInput := v3io.GetContainersInput{}
			getContainersInput.AuthenticationToken = v3iohttp.GenerateAuthenticationToken(getCommandeer.RootCommandeer.Username, getCommandeer.RootCommandeer.Password)
			getContainersInput.AccessKey = getCommandeer.RootCommandeer.AccessKey

			response, err := getCommandeer.RootCommandeer.DataPlaneContext.GetContainersSync(&getContainersInput)
			if err != nil {
				return errors.Wrap(err, "Failed to get containers")
			}

			if err := commandeer.renderContainers(response.Output.(*v3io.GetContainersOutput).Results.Containers); err != nil {
				return errors.Wrap(err, "Failed to render")
			}

			return nil
		},
	}

	commandeer.Cmd = cmd

	return commandeer, nil
}

func (c *getContainerCommandeer) renderContainers(containerInfos []v3io.ContainerInfo) error {
	renderer := renderer.NewRenderer(os.Stdout)

	switch c.RootCommandeer.Output {
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
		return renderer.RenderYAML(containerInfos) // nolint: errcheck
	case "json":
		return renderer.RenderJSON(containerInfos) // nolint: errcheck
	}

	return nil
}

// register to factory
func init() {
	v3ctl.GetCommandeerRegistrySingleton.Register("container",
		func(getCommandeer *v3ctl.GetCommandeer) (*cobra.Command, error) {
			newGetContainerCommandeer, err := newGetContainerCommandeer(getCommandeer)
			if err != nil {
				return nil, errors.Wrap(err, "Failed to get commandeer")
			}

			return newGetContainerCommandeer.Cmd, nil
		})
}
