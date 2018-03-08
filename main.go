package main

import (
	"os"

	"github.com/rosbo/pubsubbeat/cmd"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
