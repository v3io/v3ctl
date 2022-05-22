package main

import (
	"os"

	"github.com/v3io/v3ctl/pkg/v3ctl"
	_ "github.com/v3io/v3ctl/pkg/v3ctl/container"
	_ "github.com/v3io/v3ctl/pkg/v3ctl/content"
	_ "github.com/v3io/v3ctl/pkg/v3ctl/stream"

	"github.com/nuclio/errors"
)

func main() {
	rootCommandeer, err := v3ctl.NewRootCommandeer()
	if err != nil {
		errors.PrintErrorStack(os.Stderr, err, 10)
		os.Exit(1)
	}

	if err := rootCommandeer.Execute(); err != nil {
		errors.PrintErrorStack(os.Stderr, err, 10)
		os.Exit(1)
	}

	os.Exit(0)
}
