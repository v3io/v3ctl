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
	"strconv"

	"github.com/nuclio/errors"
	"github.com/spf13/cobra"
	"github.com/v3io/v3io-go/pkg/common"
	"github.com/v3io/v3io-go/pkg/dataplane/streamconsumergroup"
)

type shardOffset struct {
	ShardID  int    `json:"shard_id"`
	MemberID string `json:"member_id,omitempty"`
	Offset   uint64 `json:"offset"`
}

type getStreamConsumerGroupOffsetsCommandeer struct {
	*getStreamConsumerGroupCommandeer
	Cmd *cobra.Command
}

func newGetStreamConsumerGroupOffsetsCommandeer(getStreamConsumerGroupCommandeer *getStreamConsumerGroupCommandeer) (*getStreamConsumerGroupOffsetsCommandeer, error) {
	commandeer := &getStreamConsumerGroupOffsetsCommandeer{
		getStreamConsumerGroupCommandeer: getStreamConsumerGroupCommandeer,
	}

	cmd := &cobra.Command{
		Use:   "offsets",
		Short: "Get the offsets of a consumer group",
		RunE: func(cmd *cobra.Command, args []string) error {

			// if we got positional arguments
			if len(args) != 1 {
				return errors.New("Stream get requires a stream path")
			}

			// initialize root
			if err := getStreamConsumerGroupCommandeer.Initialize(); err != nil {
				return errors.Wrap(err, "Failed to initialize root")
			}

			shardOffsets, err := commandeer.getShardOffets(args[0])
			if err != nil {
				return errors.Wrap(err, "Failed to get shard offsets")
			}

			var records [][]string
			for _, shardOffset := range shardOffsets {
				memberID := shardOffset.MemberID
				if memberID == "" {
					memberID = "-"
				}

				records = append(records, []string{
					strconv.Itoa(shardOffset.ShardID),
					memberID,
					strconv.FormatUint(shardOffset.Offset, 10),
				})
			}

			if err := commandeer.RootCommandeer.Render(shardOffsets,
				[]string{"Shard ID", "Member ID", "Offset"},
				records); err != nil {
				return errors.Wrap(err, "Failed to render")
			}

			return nil
		},
	}

	commandeer.Cmd = cmd

	return commandeer, nil
}

func (c *getStreamConsumerGroupOffsetsCommandeer) getShardOffets(streamPath string) ([]shardOffset, error) {
	var shardOffsets []shardOffset

	streamConsumerGroup, err := streamconsumergroup.NewStreamConsumerGroup(c.RootCommandeer.Logger,
		c.name,
		streamconsumergroup.NewConfig(),
		c.RootCommandeer.Container,
		streamPath,
		0)

	if err != nil {
		return nil, errors.Wrap(err, "Failed to create consumer group")
	}

	streamConsumerGroupState, err := streamConsumerGroup.GetState()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get consumer group state")
	}

	numShards, err := streamConsumerGroup.GetNumShards()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get number of shards")
	}

	for shardID := 0; shardID < numShards; shardID++ {
		shardSequenceNumber, err := streamConsumerGroup.GetShardSequenceNumber(shardID)
		if err != nil {
			if err != streamconsumergroup.ErrShardNotFound &&
				err != streamconsumergroup.ErrShardSequenceNumberAttributeNotFound {
				return nil, errors.Wrap(err, "Failed to read shard sequence number")
			}

			// shard isn't available yet
			continue
		}

		shardMember, err := c.getShardMember(shardID, streamConsumerGroupState)
		if err != nil {
			return nil, errors.Wrap(err, "Failed to get shard member")
		}

		shardOffsets = append(shardOffsets, shardOffset{
			ShardID:  shardID,
			MemberID: shardMember,
			Offset:   shardSequenceNumber,
		})
	}

	return shardOffsets, nil
}

func (c *getStreamConsumerGroupOffsetsCommandeer) getShardMember(shardID int,
	streamConsumerGroupState *streamconsumergroup.State) (string, error) {

	for _, sessionState := range streamConsumerGroupState.SessionStates {
		if common.IntSliceContainsInt(sessionState.Shards, shardID) {
			return sessionState.MemberID, nil
		}
	}

	return "", nil
}
