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
