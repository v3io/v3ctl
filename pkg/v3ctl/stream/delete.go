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
package stream

import (
	"fmt"
	"strings"

	"github.com/v3io/v3ctl/pkg/v3ctl"

	"github.com/nuclio/errors"
	"github.com/spf13/cobra"
	v3io "github.com/v3io/v3io-go/pkg/dataplane"
)

type deleteStreamCommandeer struct {
	*v3ctl.DeleteCommandeer
}

func newDeleteStreamCommandeer(deleteCommandeer *v3ctl.DeleteCommandeer) (*deleteStreamCommandeer, error) {
	commandeer := &deleteStreamCommandeer{
		DeleteCommandeer: deleteCommandeer,
	}

	cmd := &cobra.Command{
		Use:   "stream name",
		Short: "Delete a data stream",
		RunE: func(cmd *cobra.Command, args []string) error {

			// if we got positional arguments
			if len(args) != 1 {
				return errors.New("Stream delete requires a stream name")
			}

			streamPath := args[0]

			// must end with "/"
			if !strings.HasSuffix(streamPath, "/") {
				streamPath += "/"
			}

			// initialize root
			if err := deleteCommandeer.RootCommandeer.Initialize(); err != nil {
				return errors.Wrap(err, "Failed to initialize root")
			}

			deleteStreamInput := v3io.DeleteStreamInput{}
			deleteStreamInput.Path = streamPath

			err := deleteCommandeer.RootCommandeer.Container.DeleteStreamSync(&deleteStreamInput)
			if err != nil {
				return errors.Wrap(err, "Failed to get delete stream")
			}

			fmt.Printf("Stream %s deleted successfully\n", streamPath)
			return nil
		},
	}

	commandeer.Cmd = cmd

	return commandeer, nil
}

// register to factory
func init() {
	v3ctl.DeleteCommandeerRegistrySingleton.Register("stream",
		func(deleteCommandeer *v3ctl.DeleteCommandeer) (*cobra.Command, error) {
			newDeleteStreamCommandeer, err := newDeleteStreamCommandeer(deleteCommandeer)
			if err != nil {
				return nil, errors.Wrap(err, "Failed to delete commandeer")
			}

			return newDeleteStreamCommandeer.Cmd, nil
		})
}
