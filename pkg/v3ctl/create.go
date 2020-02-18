package v3ctl

import (
	"github.com/nuclio/errors"
	"github.com/spf13/cobra"
	"github.com/v3io/registry"
)

type CreateCommandeer struct {
	Cmd            *cobra.Command
	RootCommandeer *RootCommandeer
}

func newCreateCommandeer(rootCommandeer *RootCommandeer) (*CreateCommandeer, error) {
	commandeer := &CreateCommandeer{
		RootCommandeer: rootCommandeer,
	}

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create resources",
	}

	// iterate over registry objects
	for _, createCommandeerKind := range CreateCommandeerRegistrySingleton.GetKinds() {
		createCommandeerCreatorInterface, err := CreateCommandeerRegistrySingleton.Get(createCommandeerKind)
		if err != nil {
			return nil, errors.Wrap(err, "Failed to create commandeer")
		}

		// get the creator
		createCommandeerCreator := createCommandeerCreatorInterface.(func(createCommandeer *CreateCommandeer) (*cobra.Command, error))

		createCommandeerInstance, err := createCommandeerCreator(commandeer)
		if err != nil {
			return nil, errors.Wrap(err, "Failed to create create commandeer")
		}

		// add command
		cmd.AddCommand(createCommandeerInstance)
	}

	commandeer.Cmd = cmd

	return commandeer, nil
}

func (c *CreateCommandeer) Initialize() error {
	return c.RootCommandeer.Initialize()
}

//
// Factory registry
//

// create creates a "create" commandeer
var CreateCommandeerRegistrySingleton = registry.NewRegistry("createCommandeer")
