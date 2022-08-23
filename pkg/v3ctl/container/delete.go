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
	"strconv"

	"github.com/v3io/v3ctl/pkg/v3ctl"

	"github.com/nuclio/errors"
	"github.com/spf13/cobra"
	"github.com/v3io/v3io-go/pkg/controlplane"
	v3io "github.com/v3io/v3io-go/pkg/dataplane"
	v3iohttp "github.com/v3io/v3io-go/pkg/dataplane/http"
)

type deleteContainerCommandeer struct {
	*v3ctl.DeleteCommandeer
}

func newDeleteContainerCommandeer(deleteCommandeer *v3ctl.DeleteCommandeer) (*deleteContainerCommandeer, error) {
	commandeer := &deleteContainerCommandeer{
		DeleteCommandeer: deleteCommandeer,
	}

	cmd := &cobra.Command{
		Use:   "container name",
		Short: "Delete a data container",
		RunE: func(cmd *cobra.Command, args []string) error {
			// if we got positional arguments
			if len(args) != 1 {
				return errors.New("Container delete requires a container name")
			}

			// initialize root
			if err := deleteCommandeer.Initialize(); err != nil {
				return errors.Wrap(err, "Failed to initialize root")
			}

			controlPlaneSession, err := deleteCommandeer.RootCommandeer.GetControlPlaneSession()
			if err != nil {
				return errors.Wrap(err, "Failed to get control plane session")
			}

			deleteContainerInput := v3ioc.DeleteContainerInput{}

			// resolve the container name to an ID
			deleteContainerInput.ID, err = commandeer.getContainerID(args[0])
			if err != nil {
				return errors.Wrap(err, "Failed to get container ID")
			}

			err = controlPlaneSession.DeleteContainerSync(&deleteContainerInput)
			if err != nil {
				return errors.Wrap(err, "Failed to delete container")
			}

			fmt.Printf("Container %s deleted successfully\n", args[0])
			return nil
		},
	}

	commandeer.Cmd = cmd

	return commandeer, nil
}

func (c *deleteContainerCommandeer) getContainerID(containerNameOrID string) (string, error) {

	// get containers
	getContainersInput := v3io.GetContainersInput{}
	getContainersInput.AuthenticationToken = v3iohttp.GenerateAuthenticationToken(c.RootCommandeer.Username, c.RootCommandeer.Password)
	getContainersInput.AccessKey = c.RootCommandeer.AccessKey

	response, err := c.RootCommandeer.DataPlaneContext.GetContainersSync(&getContainersInput)
	if err != nil {
		return "", errors.Wrap(err, "Failed to get containers")
	}

	// iterate over containers and look for a container whose name == name
	for _, container := range response.Output.(*v3io.GetContainersOutput).Results.Containers {
		if container.Name == containerNameOrID {
			return strconv.Itoa(container.ID), nil
		}
	}

	// couldn't find container with this name.
	// iterate over the containers and look for a container with this ID
	for _, container := range response.Output.(*v3io.GetContainersOutput).Results.Containers {
		idString := strconv.Itoa(container.ID)

		if idString == containerNameOrID {
			return idString, nil
		}
	}

	// could not find container with this name / id
	return "", errors.Errorf("Could not find container with name or ID of %s", containerNameOrID)
}

// register to factory
func init() {
	v3ctl.DeleteCommandeerRegistrySingleton.Register("container",
		func(deleteCommandeer *v3ctl.DeleteCommandeer) (*cobra.Command, error) {
			newDeleteContainerCommandeer, err := newDeleteContainerCommandeer(deleteCommandeer)
			if err != nil {
				return nil, errors.Wrap(err, "Failed to delete commandeer")
			}

			return newDeleteContainerCommandeer.Cmd, nil
		})
}
