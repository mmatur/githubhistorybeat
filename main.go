package main

import (
	"os"

	"github.com/mmatur/githubhistorybeat/cmd"
)

var name = "githubhistorybeat"

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
