package main

import (
	"os"

	"github.com/triabokon/gotagv/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
