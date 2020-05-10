/*
Copyright 2017 The Nuclio Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v3ctl

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/v3io/version-go"
)

type versionCommandeer struct {
	Cmd            *cobra.Command
	rootCommandeer *RootCommandeer
}

func newVersionCommandeer(rootCommandeer *RootCommandeer) (*versionCommandeer, error) {
	commandeer := &versionCommandeer{
		rootCommandeer: rootCommandeer,
	}

	cmd := &cobra.Command{
		Use:     "version",
		Aliases: []string{"ver"},
		Short:   "Display the version number",
		RunE: func(cmd *cobra.Command, args []string) error {

			// get the version
			fmt.Printf("Client version: %s", version.Get().String())

			return nil
		},
	}

	commandeer.Cmd = cmd

	return commandeer, nil
}
