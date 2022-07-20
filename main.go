package main

import (
	"cli-generator-poc/internal"
	"fmt"
)

func main() {
	internal.ExtractCommands()

	for _, item := range doc.Paths {
		// Grouping via item.Tags
		fmt.Printf("\nTags: %v\n", item.Post.Tags)
		fmt.Printf("Summary: %v\n", item.Post.Summary)

		// Extensions
		if item.Post.Extensions["x-cli-name"] == nil {
			fmt.Println("WARNING: No x-cli-name annotation, skipping generation.")
			continue
		}
		var command = extStr(item.Post.Extensions["x-cli-name"])
		fmt.Printf("Command: %v\n", command)
		// TODO: Create a command hierarchy when a - seperated command name is encountered
		fmt.Printf("Command name: %v\n", slug(command))

		for _, parameter := range item.Post.Parameters {
			fmt.Printf("Parameter name: %v\n", parameter.Value.Name)
			fmt.Printf("Required: %v\n", parameter.Value.Required)
		}
		fmt.Printf("OperationID: %v\n", item.Post.OperationID)

		var templateData = internal.TemplateActionsData{
			Package:   "",
			Store:     "",
			SdkModels: nil,
			Year:      0,
			Command:   internal.Command{},
		}
		internal.Generate(templateData)
	}

	// TODO: Create sdk for each command
	// Needs a store

	// TODO: Helper, get the current prod spec file, using the sha from cloud
	// TODO: Helper, get the diff of two prod spec files and only generate commands for one's that are added/updated
	// TODO: What do do for updated one - recreate the command file?...
}
