package main

import (
	"os"

	"github.com/otaviof/edsrv/pkg/edsrv/cmd"
)

func main() {
	if err := cmd.NewRoot().Cmd().Execute(); err != nil {
		os.Exit(1)
	}
}
