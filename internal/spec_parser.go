package internal

import (
	"encoding/json"
	"fmt"
	"github.com/getkin/kin-openapi/openapi3"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"strings"
	"time"
)

const (
	specLocation = "spec/openapi.yml"
)

func getParametersIgnored() []string {
	return []string{"pretty", "envelope"}
}

func ExtractCommands() []Command {
	doc := getSpec()
	var root = Command{
		Name:     "/",
		Children: make([]Command, 0),
	}

	for _, item := range doc.Paths {
		for _, operation := range item.Operations() {
			if operation != nil {
				var cliName = operation.Extensions["x-cli-name"]
				if cliName == nil {
					fmt.Println("WARNING: No x-cli-name annotation, skipping generation.")
					return nil
				}

				var command = extStr(cliName)
				var commandParts = strings.Split(command, "-")
				var packageName = commandParts[0]
				extractOperationCommand(&root, operation, packageName, commandParts)
			}
		}
	}
	return root.Children
}

func extractSchemaAsParameters(prefix string, required bool, isArray bool, schema *openapi3.Schema) []Parameter {
	var parameters = make([]Parameter, 0)
	var titleCasePrefix = cases.Title(language.Und).String(prefix)

	if schema.Type == "object" || schema.Type == "array" {
		for name, prop := range schema.Properties {
			var required = contains(schema.Required, name)
			var titleCaseName = cases.Title(language.Und).String(name)
			parameters = append(parameters, extractSchemaAsParameters(titleCasePrefix+titleCaseName, required, schema.Type == "array", prop.Value)...)
		}
	} else {
		// TODO: Get Flag
		// TODO: Get Short
		// TODO: Get Default Value
		parameters = append(parameters, Parameter{
			Name:         titleCasePrefix,
			Usage:        schema.Description,
			Flag:         "",
			Type:         schema.Type,
			IsArray:      isArray,
			Required:     required,
			Short:        "",
			DefaultValue: "",
		})
	}

	return parameters
}

func getSpec() *openapi3.T {
	// Read the openAPI spec
	loader := openapi3.NewLoader()

	doc, err := loader.LoadFromFile(specLocation)
	if err != nil {
		panic(err)
	}

	if err = doc.Validate(loader.Context); err != nil {
		panic(err)
	}

	return doc
}

func extractOperationCommand(parent *Command, operation *openapi3.Operation, packageName string, commandParts []string) {
	if len(commandParts) == 0 {
		return
	}
	var foundMatch = false

	if len(commandParts) > 1 {
		var cmdName = commandParts[0]
		for _, command := range parent.Children {
			if command.Name == cmdName {
				foundMatch = true
				break
			}
		}
		if !foundMatch {
			var newCommand = Command{
				Year:     time.Now().Year(),
				Package:  packageName,
				Name:     commandParts[0],
				Children: make([]Command, 0),
			}
			parent.Children = append(parent.Children, newCommand)
		}
		extractOperationCommand(&parent.Children[len(parent.Children)-1], operation, packageName, commandParts[1:])
		return
	}

	var leafCommand = Command{
		Year:           time.Now().Year(),
		Name:           commandParts[len(commandParts)-1],
		Short:          "TODO",
		Long:           "TODO",
		Parameters:     make([]Parameter, 0),
		Header:         "TODO",
		TemplateOutput: "TODO",
	}

	for _, parameter := range operation.Parameters {
		if contains(getParametersIgnored(), parameter.Value.Name) {
			continue
		}
		short := formatExtValue(parameter, "x-cli-short", "")
		defaultValue := formatExtValue(parameter, "x-cli-default", "")
		flag := formatExtValue(parameter, "x-cli-flag", "<MISSING_FLAG>")
		cliDescription := formatExtValue(parameter, "x-cli-description", parameter.Value.Description)
		leafCommand.Parameters = append(leafCommand.Parameters, Parameter{
			Name:         parameter.Value.Name,
			Usage:        cliDescription,
			Flag:         flag,
			Type:         parameter.Value.Schema.Value.Type,
			Required:     parameter.Value.Required,
			Short:        short,
			DefaultValue: defaultValue,
		})
	}
	schema := operation.RequestBody.Value.Content.Get("application/json").Schema.Value
	leafCommand.Parameters = append(leafCommand.Parameters, extractSchemaAsParameters("", false, false, schema)...)

	parent.Children = append(parent.Children, leafCommand)
	return
}

func formatExtValue(parameter *openapi3.ParameterRef, extKey string, defaultValue string) string {
	shortRaw := parameter.Value.Extensions[extKey]
	if shortRaw != nil {
		return extStr(shortRaw)
	}
	return defaultValue
}

// extStr returns the string value of an OpenAPI extension stored as a JSON
// raw message.
func extStr(i interface{}) (decoded string) {
	if err := json.Unmarshal(i.(json.RawMessage), &decoded); err != nil {
		panic(err)
	}

	return
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
