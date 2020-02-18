package container

import (
	"strconv"

	"github.com/v3io/v3ctl/pkg/v3ctl"

	"github.com/nuclio/errors"
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

			defer response.Release()

			containerInfos := response.Output.(*v3io.GetContainersOutput).Results.Containers

			var records [][]string
			for _, containerInfo := range containerInfos {
				records = append(records, []string{
					strconv.Itoa(containerInfo.ID),
					containerInfo.Name,
					containerInfo.CreationDate,
				})
			}

			if err := commandeer.RootCommandeer.Render(containerInfos,
				[]string{"ID", "Name", "Creation date"},
				records); err != nil {
				return errors.Wrap(err, "Failed to render")
			}

			return nil
		},
	}

	commandeer.Cmd = cmd

	return commandeer, nil
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
