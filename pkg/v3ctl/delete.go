package v3ctl

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/nuclio/errors"
	"github.com/spf13/cobra"
	"github.com/v3io/v3io-go/pkg/controlplane"
	"github.com/v3io/v3io-go/pkg/dataplane"
	"github.com/v3io/v3io-go/pkg/dataplane/http"
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

			// resolve the container name to an ID
			deleteContainerInput.ID, err = deleteCommandeer.getContainerID(args[0])
			if err != nil {
				return errors.Wrap(err, "Failed to get container ID")
			}

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

func (c *deleteCommandeer) getContainerID(containerNameOrID string) (string, error) {

	// get containers
	getContainersInput := v3io.GetContainersInput{}
	getContainersInput.AuthenticationToken = v3iohttp.GenerateAuthenticationToken(c.rootCommandeer.username, c.rootCommandeer.password)
	getContainersInput.AccessKey = c.rootCommandeer.accessKey

	response, err := c.rootCommandeer.dataPlaneContext.GetContainersSync(&getContainersInput)
	if err != nil {
		return "", errors.Wrap(err, "Failed to get containers")
	}

	// iterate over containers and look for a container whose name == name
	for _, container := range response.Output.(*v3io.GetContainersOutput).Results.Containers {
		if container.Name == containerNameOrID {
			return strconv.Itoa(container.ID), nil
		}
	}

	// couldn't find container with this name.
	// iterate over the containers and look for a container with this ID
	for _, container := range response.Output.(*v3io.GetContainersOutput).Results.Containers {
		idString := strconv.Itoa(container.ID)

		if idString == containerNameOrID {
			return idString, nil
		}
	}

	// could not find container with this name / id
	return "", errors.Errorf("Could not find container with name or ID of %s", containerNameOrID)
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
