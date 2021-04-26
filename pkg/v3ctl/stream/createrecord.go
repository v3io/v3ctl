package stream

import (
	"fmt"
	"strings"

	"github.com/nuclio/errors"
	"github.com/spf13/cobra"
	v3io "github.com/v3io/v3io-go/pkg/dataplane"
)

type createRecordCommandeer struct {
	*createStreamCommandeer
	shardID      int
	clientInfo   string
	data         string
	partitionKey string
}

func newCreateRecordCommandeer(createStreamCommandeer *createStreamCommandeer) (*createRecordCommandeer, error) {
	commandeer := &createRecordCommandeer{
		createStreamCommandeer: createStreamCommandeer,
	}

	cmd := &cobra.Command{
		Use:   "record",
		Short: "Create records in a stream",
		RunE: func(cmd *cobra.Command, args []string) error {

			// if we got positional arguments
			if len(args) != 1 {
				return errors.New("Stream create requires a stream path")
			}

			streamPath := args[0]

			// must end with "/"
			if !strings.HasSuffix(streamPath, "/") {
				streamPath += "/"
			}

			// initialize root
			if err := createStreamCommandeer.Initialize(); err != nil {
				return errors.Wrap(err, "Failed to initialize root")
			}

			response, err := commandeer.RootCommandeer.Container.PutRecordsSync(&v3io.PutRecordsInput{
				Path: streamPath,
				Records: []*v3io.StreamRecord{
					{
						ShardID:      &commandeer.shardID,
						Data:         []byte(commandeer.data),
						ClientInfo:   []byte(commandeer.clientInfo),
						PartitionKey: commandeer.partitionKey,
					},
				},
			})

			if err != nil {
				return errors.Wrap(err, "Failed to get create stream")
			}

			putRecordsResponse, ok := response.Output.(v3io.PutRecordsOutput)
			if ok && putRecordsResponse.FailedRecordCount != 0 {
				return errors.Errorf("Failed to put all records, FailedRecordCount: %d",
					putRecordsResponse.FailedRecordCount)
			}

			defer response.Release()

			fmt.Printf("Wrote %d bytes to %s:%d\n",
				len(commandeer.data),
				streamPath,
				commandeer.shardID)

			return nil
		},
	}

	cmd.Flags().IntVar(&commandeer.shardID, "shard-id", 1, "")
	cmd.Flags().StringVar(&commandeer.clientInfo, "client-info", "", "")
	cmd.Flags().StringVar(&commandeer.partitionKey, "partition-key", "", "")
	cmd.Flags().StringVar(&commandeer.data, "data", "", "")

	commandeer.Cmd = cmd

	return commandeer, nil
}
