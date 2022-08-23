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

	"github.com/nuclio/errors"
	"github.com/spf13/cobra"
	"github.com/v3io/v3io-go/pkg/dataplane/streamconsumergroup"
)

type getStreamConsumerGroupCommandeer struct {
	*getStreamCommandeer
	name string
}

func newGetStreamConsumerGroupCommandeer(getStreamCommandeer *getStreamCommandeer) (*getStreamConsumerGroupCommandeer, error) {
	commandeer := &getStreamConsumerGroupCommandeer{
		getStreamCommandeer: getStreamCommandeer,
	}

	cmd := &cobra.Command{
		Use:   "consumer-group",
		Short: "Get consumer group related information",
		RunE: func(cmd *cobra.Command, args []string) error {

			// if we got positional arguments
			if len(args) != 1 {
				return errors.New("Stream get requires a stream path")
			}

			// initialize root
			if err := getStreamCommandeer.Initialize(); err != nil {
				return errors.Wrap(err, "Failed to initialize root")
			}

			path := args[0]

			streamConsumerGroup, err := streamconsumergroup.NewStreamConsumerGroup(getStreamCommandeer.RootCommandeer.Logger,
				commandeer.name,
				streamconsumergroup.NewConfig(),
				getStreamCommandeer.RootCommandeer.Container,
				path,
				0)

			if err != nil {
				return errors.Wrap(err, "Failed to create consumer group")
			}

			streamConsumerGroupState, err := streamConsumerGroup.GetState()
			if err != nil {
				return errors.Wrap(err, "Failed to get consumer group state")
			}

			fmt.Println(streamConsumerGroupState.SessionStates)

			return nil
		},
	}

	cmd.PersistentFlags().StringVar(&commandeer.name, "consumer-group-name", "", "")

	getStreamConsumerGroupMembers, err := newGetStreamConsumerGroupMembersCommandeer(commandeer)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create create record")
	}

	getStreamConsumerGroupOffsets, err := newGetStreamConsumerGroupOffsetsCommandeer(commandeer)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create create record")
	}

	cmd.AddCommand(getStreamConsumerGroupMembers.Cmd)
	cmd.AddCommand(getStreamConsumerGroupOffsets.Cmd)

	commandeer.Cmd = cmd

	return commandeer, nil
}
