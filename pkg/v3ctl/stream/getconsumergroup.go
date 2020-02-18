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
