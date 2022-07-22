package main

import (
	"cli-generator-poc/internal"
)

const (
	generatedPath = "generated"
)

func main() {
	var commands = internal.ExtractCommands()
	for _, command := range commands {
		// TODO: Create sdk for each command
		// Needs a store

		// Generate the command template
		internal.Generate(generatedPath, command)
	}

	// TODO: Helper, get the current prod spec file, using the sha from cloud
	// TODO: Helper, get the diff of two prod spec files and only generate commands for one's that are added/updated
	// TODO: What do do for updated one - recreate the command file?...
}
