package v3ctl

import (
	"github.com/nuclio/errors"
	"github.com/spf13/cobra"
	"github.com/v3io/registry"
)

type DeleteCommandeer struct {
	Cmd            *cobra.Command
	RootCommandeer *RootCommandeer
}

func newDeleteCommandeer(rootCommandeer *RootCommandeer) (*DeleteCommandeer, error) {
	commandeer := &DeleteCommandeer{
		RootCommandeer: rootCommandeer,
	}

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete resources",
	}

	// iterate over registry objects
	for _, deleteCommandeerKind := range DeleteCommandeerRegistrySingleton.GetKinds() {
		deleteCommandeerCreatorInterface, err := DeleteCommandeerRegistrySingleton.Get(deleteCommandeerKind)
		if err != nil {
			return nil, errors.Wrap(err, "Failed to delete commandeer")
		}

		// get the creator
		deleteCommandeerCreator := deleteCommandeerCreatorInterface.(func(deleteCommandeer *DeleteCommandeer) (*cobra.Command, error))

		deleteCommandeerInstance, err := deleteCommandeerCreator(commandeer)
		if err != nil {
			return nil, errors.Wrap(err, "Failed to delete delete commandeer")
		}

		// add command
		cmd.AddCommand(deleteCommandeerInstance)
	}

	commandeer.Cmd = cmd

	return commandeer, nil
}

func (c *DeleteCommandeer) Initialize() error {
	return c.RootCommandeer.Initialize()
}

//
// Factory registry
//

// delete deletes a "delete" commandeer
var DeleteCommandeerRegistrySingleton = registry.NewRegistry("deleteCommandeer")
