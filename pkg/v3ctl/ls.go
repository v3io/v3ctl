package v3ctl

import (
	"encoding/json"
	"fmt"
	"github.com/nuclio/errors"
	"github.com/spf13/cobra"
	v3io "github.com/v3io/v3io-go/pkg/dataplane"
)

type lsCommandeer struct {
	rootCommandeer *RootCommandeer
	cmd            *cobra.Command
	withAttributes bool
}

func newLsCommandeer(rootCommandeer *RootCommandeer) *lsCommandeer {
	commandeer := &lsCommandeer{
		rootCommandeer: rootCommandeer,
	}

	cmd := &cobra.Command{
		Use:   "ls <dir> [-l]",
		Short: "List files",
		RunE: func(cmd *cobra.Command, args []string) error {

			if len(args) != 1 {
				return errors.New("ls requires a directory")
			}

			dir := args[0]

			// initialize root
			if err := commandeer.rootCommandeer.initialize(); err != nil {
				return errors.Wrap(err, "Failed to initialize root")
			}

			input := v3io.GetContainerContentsInput{Path: dir, GetAllAttributes: commandeer.withAttributes}
			for {
				res, err := rootCommandeer.container.GetContainerContentsSync(&input)
				if err != nil {
					return err
				}
				out := res.Output.(*v3io.GetContainerContentsOutput)
				err = printEntries(out.CommonPrefixes, out.Contents, commandeer.withAttributes)
				if err != nil {
					return err
				}
				if !out.IsTruncated {
					break
				}
				input.Marker = out.NextMarker
			}

			return nil
		},
	}

	cmd.Flags().BoolVarP(&commandeer.withAttributes, "with-attributes", "l", false, "Show attributes")

	commandeer.cmd = cmd

	return commandeer
}

func printEntries(commonPrefixes []v3io.CommonPrefix, contents []v3io.Content, withDetail bool) error {
	if withDetail {
		for _, commonPrefix := range commonPrefixes {
			line, err := json.Marshal(commonPrefix)
			if err != nil {
				return errors.Wrapf(err, "failed to marshal common prefix %s", commonPrefix.Prefix)
			}
			fmt.Printf("%s\n", line)
		}
		for _, content := range contents {
			line, err := json.Marshal(content)
			if err != nil {
				return errors.Wrapf(err, "failed to marshal content %s", content.Key)
			}
			fmt.Printf("%s\n", line)
		}
	} else {
		for _, commonPrefix := range commonPrefixes {
			fmt.Println(commonPrefix.Prefix)
		}
		for _, content := range contents {
			fmt.Println(content.Key)
		}
	}
	return nil
}
