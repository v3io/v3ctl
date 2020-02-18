package v3ctl

import (
	"github.com/nuclio/errors"
	"github.com/spf13/cobra"
	"github.com/v3io/registry"
)

type GetCommandeer struct {
	Cmd            *cobra.Command
	RootCommandeer *RootCommandeer
}

func newGetCommandeer(rootCommandeer *RootCommandeer) (*GetCommandeer, error) {
	commandeer := &GetCommandeer{
		RootCommandeer: rootCommandeer,
	}

	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get resources",
	}

	// iterate over registry objects
	for _, getCommandeerKind := range GetCommandeerRegistrySingleton.GetKinds() {
		getCommandeerCreatorInterface, err := GetCommandeerRegistrySingleton.Get(getCommandeerKind)
		if err != nil {
			return nil, errors.Wrap(err, "Failed to get commandeer")
		}

		// get the creator
		getCommandeerCreator := getCommandeerCreatorInterface.(func(getCommandeer *GetCommandeer) (*cobra.Command, error))

		getCommandeerInstance, err := getCommandeerCreator(commandeer)
		if err != nil {
			return nil, errors.Wrap(err, "Failed to get get commandeer")
		}

		// add command
		cmd.AddCommand(getCommandeerInstance)
	}

	commandeer.Cmd = cmd

	return commandeer, nil
}

func (c *GetCommandeer) Initialize() error {
	return c.RootCommandeer.Initialize()
}

//
// Factory registry
//

// get gets a "get" commandeer
var GetCommandeerRegistrySingleton = registry.NewRegistry("getCommandeer")
