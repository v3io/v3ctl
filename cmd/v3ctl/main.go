package main

import (
	"os"

	"github.com/iguazio/v3ctl/pkg/v3ctl"

	"github.com/nuclio/errors"
)

func main() {
	if err := v3ctl.NewRootCommandeer().Execute(); err != nil {
		errors.PrintErrorStack(os.Stderr, err, 5)
		os.Exit(1)
	}

	os.Exit(0)
}
