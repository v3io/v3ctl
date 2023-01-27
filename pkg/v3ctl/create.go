/*
Copyright 2019 Iguazio Systems Ltd.

Licensed under the Apache License, Version 2.0 (the "License") with
an addition restriction as set forth herein. You may not use this
file except in compliance with the License. You may obtain a copy of
the License at http://www.apache.org/licenses/LICENSE-2.0.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
implied. See the License for the specific language governing
permissions and limitations under the License.

In addition, you may not use the software for any purposes that are
illegal under applicable law, and the grant of the foregoing license
under the Apache 2.0 license is conditioned upon your compliance with
such restriction.
*/

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
