package stream

import (
	"fmt"

	"github.com/v3io/v3ctl/pkg/v3ctl"

	"github.com/nuclio/errors"
	"github.com/spf13/cobra"
	v3io "github.com/v3io/v3io-go/pkg/dataplane"
	v3iohttp "github.com/v3io/v3io-go/pkg/dataplane/http"
)

type getStreamCommandeer struct {
	*v3ctl.GetCommandeer
}

func newGetStreamCommandeer(getCommandeer *v3ctl.GetCommandeer) (*getStreamCommandeer, error) {
	commandeer := &getStreamCommandeer{
		GetCommandeer: getCommandeer,
	}

	cmd := &cobra.Command{
		Use:   "stream name",
		Short: "Get a data stream",
		RunE: func(cmd *cobra.Command, args []string) error {
			// if we got positional arguments
			if len(args) != 1 {
				return errors.New("Stream get requires a stream path")
			}

			// initialize root
			if err := getCommandeer.Initialize(); err != nil {
				return errors.Wrap(err, "Failed to initialize root")
			}

			path := args[0]

			// populate request
			getItemsInput := &v3io.GetItemsInput{}
			getItemsInput.Path = path
			getItemsInput.ContainerName = getCommandeer.RootCommandeer.ContainerName
			getItemsInput.AuthenticationToken = v3iohttp.GenerateAuthenticationToken(getCommandeer.RootCommandeer.Username, getCommandeer.RootCommandeer.Password)
			getItemsInput.AccessKey = getCommandeer.RootCommandeer.AccessKey

			response, err := getCommandeer.RootCommandeer.DataPlaneContext.GetItemsSync(getItemsInput)

			if err != nil {
				return errors.Wrapf(err, "Failed to get container contents at %s", args[0])
			}

			for _, content := range response.Output.(*v3io.GetItemsOutput).Items {
				fmt.Println(content)
			}

			return nil
		},
	}

	commandeer.Cmd = cmd

	return commandeer, nil
}

// register to factory
func init() {
	v3ctl.GetCommandeerRegistrySingleton.Register("stream",
		func(getCommandeer *v3ctl.GetCommandeer) (*cobra.Command, error) {
			newGetStreamCommandeer, err := newGetStreamCommandeer(getCommandeer)
			if err != nil {
				return nil, errors.Wrap(err, "Failed to get commandeer")
			}

			return newGetStreamCommandeer.Cmd, nil
		})
}
