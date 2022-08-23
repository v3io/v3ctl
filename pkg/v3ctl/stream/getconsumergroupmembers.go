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
	v3ioerrors "github.com/v3io/v3io-go/pkg/errors"
)

type getStreamConsumerGroupMembersCommandeer struct {
	*getStreamConsumerGroupCommandeer
	Cmd *cobra.Command
}

func newGetStreamConsumerGroupMembersCommandeer(getStreamConsumerGroupCommandeer *getStreamConsumerGroupCommandeer) (*getStreamConsumerGroupMembersCommandeer, error) {
	commandeer := &getStreamConsumerGroupMembersCommandeer{
		getStreamConsumerGroupCommandeer: getStreamConsumerGroupCommandeer,
	}

	cmd := &cobra.Command{
		Use:   "members",
		Short: "Get the members of a consumer group",
		RunE: func(cmd *cobra.Command, args []string) error {

			// if we got positional arguments
			if len(args) != 1 {
				return errors.New("Stream get requires a stream path")
			}

			// initialize root
			if err := getStreamConsumerGroupCommandeer.Initialize(); err != nil {
				return errors.Wrap(err, "Failed to initialize root")
			}

			path := args[0]

			streamConsumerGroup, err := streamconsumergroup.NewStreamConsumerGroup(getStreamConsumerGroupCommandeer.RootCommandeer.Logger,
				commandeer.name,
				streamconsumergroup.NewConfig(),
				getStreamConsumerGroupCommandeer.RootCommandeer.Container,
				path,
				0)

			if err != nil {
				return errors.Wrap(err, "Failed to create consumer group")
			}

			streamConsumerGroupState, err := streamConsumerGroup.GetState()
			if err != nil {
				if errors.Cause(err) == v3ioerrors.ErrNotFound {
					if err := commandeer.RootCommandeer.Render([]string{},
						[]string{"Member ID", "Last Heartbeat", "Shards"},
						[][]string{}); err != nil {
						return errors.Wrap(err, "Failed to render")
					}

					return nil
				}

				return errors.Wrap(err, "Failed to get consumer group state")
			}

			var records [][]string
			for _, sessionState := range streamConsumerGroupState.SessionStates {
				records = append(records, []string{
					sessionState.MemberID,
					sessionState.LastHeartbeat.UTC().Format("02/01/2006 15:04:05.000"),
					fmt.Sprintf("%d: %v", len(sessionState.Shards), sessionState.Shards),
				})
			}

			if err := commandeer.RootCommandeer.Render(streamConsumerGroupState.SessionStates,
				[]string{"Member ID", "Last Heartbeat", "Shards"},
				records); err != nil {
				return errors.Wrap(err, "Failed to render")
			}

			return nil
		},
	}

	commandeer.Cmd = cmd

	return commandeer, nil
}
