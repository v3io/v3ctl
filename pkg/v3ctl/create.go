package v3ctl

import (
	"fmt"
	"strings"

	"github.com/v3io/v3io-go/pkg/controlplane"
	"github.com/v3io/v3io-go/pkg/dataplane"

	"github.com/nuclio/errors"
	"github.com/spf13/cobra"
)

type createCommandeer struct {
	cmd            *cobra.Command
	rootCommandeer *RootCommandeer
}

func newCreateCommandeer(rootCommandeer *RootCommandeer) *createCommandeer {
	commandeer := &createCommandeer{
		rootCommandeer: rootCommandeer,
	}

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create resources",
	}

	createContainerCommand := newCreateContainerCommandeer(commandeer).cmd
	createStreamCommand := newCreateStreamCommandeer(commandeer).cmd

	cmd.AddCommand(
		createContainerCommand,
		createStreamCommand,
	)

	commandeer.cmd = cmd

	return commandeer
}

type createContainerCommandeer struct {
	*createCommandeer
}

func newCreateContainerCommandeer(createCommandeer *createCommandeer) *createContainerCommandeer {
	commandeer := &createContainerCommandeer{
		createCommandeer: createCommandeer,
	}

	cmd := &cobra.Command{
		Use:   "container name",
		Short: "Create a data container",
		RunE: func(cmd *cobra.Command, args []string) error {

			// if we got positional arguments
			if len(args) != 1 {
				return errors.New("Container create requires a container name")
			}

			// initialize root
			if err := createCommandeer.rootCommandeer.initialize(); err != nil {
				return errors.Wrap(err, "Failed to initialize root")
			}

			controlPlaneSession, err := createCommandeer.rootCommandeer.getControlPlaneSession()
			if err != nil {
				return errors.Wrap(err, "Failed to get control plane session")
			}

			createContainerInput := v3ioc.CreateContainerInput{}
			createContainerInput.Name = args[0]

			_, err = controlPlaneSession.CreateContainerSync(&createContainerInput)
			if err != nil {
				return errors.Wrap(err, "Failed to create container")
			}

			fmt.Printf("Container %s created successfully\n", args[0])
			return nil
		},
	}

	commandeer.cmd = cmd

	return commandeer
}

type createStreamCommandeer struct {
	*createCommandeer
	createStreamInput v3io.CreateStreamInput
}

func newCreateStreamCommandeer(createCommandeer *createCommandeer) *createStreamCommandeer {
	commandeer := &createStreamCommandeer{
		createCommandeer: createCommandeer,
	}

	cmd := &cobra.Command{
		Use:   "stream name [flags]",
		Short: "Create a stream",
		RunE: func(cmd *cobra.Command, args []string) error {

			// if we got positional arguments
			if len(args) != 1 {
				return errors.New("Stream create requires a stream name")
			}

			streamPath := args[0]

			// must end with "/"
			if !strings.HasSuffix(streamPath, "/") {
				streamPath += "/"
			}

			// initialize root
			if err := createCommandeer.rootCommandeer.initialize(); err != nil {
				return errors.Wrap(err, "Failed to initialize root")
			}

			commandeer.createStreamInput.Path = streamPath

			err := createCommandeer.rootCommandeer.container.CreateStreamSync(&commandeer.createStreamInput)
			if err != nil {
				return errors.Wrap(err, "Failed to get create stream")
			}

			fmt.Printf("Stream %s created successfully\n", streamPath)
			return nil
		},
	}

	cmd.Flags().IntVarP(&commandeer.createStreamInput.ShardCount, "shards", "", 1, "Number of shards")
	cmd.Flags().IntVarP(&commandeer.createStreamInput.RetentionPeriodHours, "retention-period-hours", "", 1, "Data retention period, in hours")

	commandeer.cmd = cmd

	return commandeer
}
