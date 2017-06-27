package main

import (
	"os"

	"github.com/elastic/beats/libbeat/beat"

	"github.com/mmatur/githubhistorybeat/beater"
)

var version = "1.0.0-SNAPSHOT"
var name = "githubhistorybeat"

func main() {
	err := beat.Run(name, version, beater.New)
	if err != nil {
		os.Exit(1)
	}
}
