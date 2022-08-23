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
