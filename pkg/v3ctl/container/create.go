package container

import (
	"fmt"

	"github.com/v3io/v3ctl/pkg/v3ctl"

	"github.com/nuclio/errors"
	"github.com/spf13/cobra"
	"github.com/v3io/v3io-go/pkg/controlplane"
)

type createContainerCommandeer struct {
	*v3ctl.CreateCommandeer
}

func newCreateContainerCommandeer(createCommandeer *v3ctl.CreateCommandeer) (*createContainerCommandeer, error) {
	commandeer := &createContainerCommandeer{
		CreateCommandeer: createCommandeer,
	}

	cmd := &cobra.Command{
		Use:   "container name",
		Short: "Create a data container",
		RunE: func(cmd *cobra.Command, args []string) error {

			// if we got positional arguments
			if len(args) != 1 {
				return errors.New("Container create requires a container name")
			}

			// initialize root
			if err := createCommandeer.Initialize(); err != nil {
				return errors.Wrap(err, "Failed to initialize root")
			}

			controlPlaneSession, err := createCommandeer.RootCommandeer.GetControlPlaneSession()
			if err != nil {
				return errors.Wrap(err, "Failed to get control plane session")
			}

			createContainerInput := v3ioc.CreateContainerInput{}
			createContainerInput.Name = args[0]

			_, err = controlPlaneSession.CreateContainerSync(&createContainerInput)
			if err != nil {
				return errors.Wrap(err, "Failed to create container")
			}

			fmt.Printf("Container %s created successfully\n", args[0])
			return nil
		},
	}

	commandeer.Cmd = cmd

	return commandeer, nil
}

// register to factory
func init() {
	v3ctl.CreateCommandeerRegistrySingleton.Register("container",
		func(createCommandeer *v3ctl.CreateCommandeer) (*cobra.Command, error) {
			newCreateContainerCommandeer, err := newCreateContainerCommandeer(createCommandeer)
			if err != nil {
				return nil, errors.Wrap(err, "Failed to create commandeer")
			}

			return newCreateContainerCommandeer.Cmd, nil
		})
}
