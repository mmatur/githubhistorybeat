package main

import (
	"os"

	"github.com/elastic/beats/libbeat/beat"

	"github.com/mmatur/githubhistorybeat/beater"
)

var name = "githubhistorybeat"

func main() {
	err := beat.Run(name, "", beater.New)
	if err != nil {
		os.Exit(1)
	}
}
