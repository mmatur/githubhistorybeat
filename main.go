package main

import (
	"os"

	"github.com/elastic/beats/libbeat/beat"

	"github.com/mmatur/githubhistorybeat/beater"
)

func main() {
	err := beat.Run("githubhistorybeat", "", beater.New)
	if err != nil {
		os.Exit(1)
	}
}
