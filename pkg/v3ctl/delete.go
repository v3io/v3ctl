package v3ctl

import (
	"fmt"
	"strings"

	"github.com/v3io/v3io-go/pkg/controlplane"
	"github.com/v3io/v3io-go/pkg/dataplane"

	"github.com/nuclio/errors"
	"github.com/spf13/cobra"
)

type deleteCommandeer struct {
	cmd            *cobra.Command
	rootCommandeer *RootCommandeer
}

func newDeleteCommandeer(rootCommandeer *RootCommandeer) *deleteCommandeer {
	commandeer := &deleteCommandeer{
		rootCommandeer: rootCommandeer,
	}

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete resources",
	}

	deleteContainerCommand := newDeleteContainerCommandeer(commandeer).cmd
	deleteStreamCommand := newDeleteStreamCommandeer(commandeer).cmd

	cmd.AddCommand(
		deleteContainerCommand,
		deleteStreamCommand,
	)

	commandeer.cmd = cmd

	return commandeer
}

type deleteContainerCommandeer struct {
	*deleteCommandeer
}

func newDeleteContainerCommandeer(deleteCommandeer *deleteCommandeer) *deleteContainerCommandeer {
	commandeer := &deleteContainerCommandeer{
		deleteCommandeer: deleteCommandeer,
	}

	cmd := &cobra.Command{
		Use:   "container name",
		Short: "Delete a data container",
		RunE: func(cmd *cobra.Command, args []string) error {

			// if we got positional arguments
			if len(args) != 1 {
				return errors.New("Container delete requires a container name")
			}

			// initialize root
			if err := deleteCommandeer.rootCommandeer.initialize(); err != nil {
				return errors.Wrap(err, "Failed to initialize root")
			}

			controlPlaneSession, err := deleteCommandeer.rootCommandeer.getControlPlaneSession()
			if err != nil {
				return errors.Wrap(err, "Failed to get control plane session")
			}

			deleteContainerInput := v3ioc.DeleteContainerInput{}
			deleteContainerInput.ID = args[0]

			err = controlPlaneSession.DeleteContainerSync(&deleteContainerInput)
			if err != nil {
				return errors.Wrap(err, "Failed to delete container")
			}

			fmt.Printf("Container %s deleted successfully\n", args[0])
			return nil
		},
	}

	commandeer.cmd = cmd

	return commandeer
}

type deleteStreamCommandeer struct {
	*deleteCommandeer
}

func newDeleteStreamCommandeer(deleteCommandeer *deleteCommandeer) *deleteStreamCommandeer {
	commandeer := &deleteStreamCommandeer{
		deleteCommandeer: deleteCommandeer,
	}

	cmd := &cobra.Command{
		Use:   "stream name [flags]",
		Short: "Delete a stream",
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
			if err := deleteCommandeer.rootCommandeer.initialize(); err != nil {
				return errors.Wrap(err, "Failed to initialize root")
			}

			deleteStreamInput := v3io.DeleteStreamInput{}
			deleteStreamInput.Path = streamPath

			err := deleteCommandeer.rootCommandeer.container.DeleteStreamSync(&deleteStreamInput)
			if err != nil {
				return errors.Wrap(err, "Failed to get delete stream")
			}

			fmt.Printf("Stream %s deleted successfully\n", streamPath)
			return nil
		},
	}

	commandeer.cmd = cmd

	return commandeer
}
