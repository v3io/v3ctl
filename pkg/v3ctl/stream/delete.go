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
