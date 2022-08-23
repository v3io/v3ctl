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
package content

import (
	"fmt"
	"github.com/nuclio/errors"
	"github.com/spf13/cobra"
	"github.com/v3io/v3ctl/pkg/v3ctl"
	v3io "github.com/v3io/v3io-go/pkg/dataplane"
	"time"
)

type getContentCommandeer struct {
	*v3ctl.GetCommandeer
	dirsOnly         bool
	getAllAttributes bool
}

func newGetContentCommandeer(getCommandeer *v3ctl.GetCommandeer) (*getContentCommandeer, error) {

	commandeer := &getContentCommandeer{
		GetCommandeer: getCommandeer,
	}

	cmd := &cobra.Command{
		Use:   "content pathname",
		Short: "Get a prefix content",
		RunE: func(cmd *cobra.Command, args []string) error {

			// if we got positional arguments
			if len(args) != 1 {
				return errors.New("Content get requires a pathname")
			}

			// initialize root
			if err := getCommandeer.Initialize(); err != nil {
				return errors.Wrap(err, "Failed to initialize root")
			}

			getContainerContentsInput := &v3io.GetContainerContentsInput{
				GetAllAttributes: commandeer.getAllAttributes,
				DirectoriesOnly:  commandeer.dirsOnly,
				Limit:            1000,
				DataPlaneInput: v3io.DataPlaneInput{
					URL:                    getCommandeer.RootCommandeer.WebapiURL,
					AccessKey:              getCommandeer.RootCommandeer.AccessKey,
					Timeout:                time.Duration(60) * time.Second,
					IncludeResponseInError: true,
				},
			}

			// Get subdirectories in path
			getContainerContentsInput.Path = args[0]
			getContainerContentsInput.ContainerName = getCommandeer.RootCommandeer.ContainerName
			for {

				response, err := getCommandeer.RootCommandeer.DataPlaneContext.GetContainerContentsSync(getContainerContentsInput)

				if err != nil {
					return errors.Wrapf(err, "Failed to get prefix contents at %s", args[0])
				}

				defer response.Release()

				getContainerContentsOutput := response.Output.(*v3io.GetContainerContentsOutput)

				fmt.Println("prefixes:")
				for _, prefix := range getContainerContentsOutput.CommonPrefixes {
					fmt.Printf("%+v\n", prefix)
				}

				fmt.Println("files:")
				for _, object := range getContainerContentsOutput.Contents {
					fmt.Printf("%+v\n", object)
				}

				if !getContainerContentsOutput.IsTruncated || len(getContainerContentsOutput.NextMarker) == 0 {
					getContainerContentsInput.Marker = ""
					break
				}
				getContainerContentsInput.Marker = getContainerContentsOutput.NextMarker

			}

			return nil
		},
	}

	cmd.Flags().BoolVar(&commandeer.dirsOnly, "dirs-only", false, "retrieve directories only")
	cmd.Flags().BoolVar(&commandeer.getAllAttributes, "get-all-attrs", true, "retrieve all directory attributes")

	commandeer.Cmd = cmd

	return commandeer, nil
}

// register to factory
func init() {
	v3ctl.GetCommandeerRegistrySingleton.Register("content",
		func(getCommandeer *v3ctl.GetCommandeer) (*cobra.Command, error) {
			newGetContentCommandeer, err := newGetContentCommandeer(getCommandeer)
			if err != nil {
				return nil, errors.Wrap(err, "Failed to get commandeer")
			}

			return newGetContentCommandeer.Cmd, nil
		})
}
