package main

import (
	"os"

	"sqlbeat/cmd"
)

func main() {

	if err := cmd.RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
