package internal

import (
	"encoding/json"
	"fmt"
	"github.com/getkin/kin-openapi/openapi3"
	"time"
)

const (
	specLocation = "spec/openapi.yml"
)

func ExtractCommands() []TemplateActionsData {
	doc := validateSpec()
	commands := make([]TemplateActionsData, len(doc.Paths))

	for _, item := range doc.Paths {
		for _, operation := range item.Operations() {
			if operation != nil {
				command := extractCommand(operation)
				if command != nil {
					commands = append(commands, command)
				}
			}
		}
	}

	return commands
}

func validateSpec() *openapi3.T {
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

func extractCommand(operation *openapi3.Operation) *TemplateActionsData {
	if operation.Extensions["x-cli-name"] == nil {
		fmt.Println("WARNING: No x-cli-name annotation, skipping generation.")
		return nil
	}
	var command = extStr(operation.Extensions["x-cli-name"])

	var tmplData = &TemplateActionsData{Year: time.Now().Year()}

	tmplData.Package = operation.OperationID
	tmplData.Command = Command{
		Name:           command,
		Short:          "TODO",
		Long:           "TODO",
		Parameters:     make([]Parameter, len(operation.Parameters)),
		Header:         "TODO",
		TemplateOutput: "TODO",
	}

	for _, parameter := range operation.Parameters {
		tmplData.Command.Parameters = append(tmplData.Command.Parameters, Parameter{
			Name:     parameter.Value.Name,
			Type:     parameter.Value.Schema.Value.Type,
			Required: parameter.Value.Required,
		})
	}
	return tmplData
}

// extStr returns the string value of an OpenAPI extension stored as a JSON
// raw message.
func extStr(i interface{}) (decoded string) {
	if err := json.Unmarshal(i.(json.RawMessage), &decoded); err != nil {
		panic(err)
	}

	return
}

//func slug(operationID string) string {
//	transformed := strings.ToLower(operationID)
//	transformed = strings.Replace(transformed, "_", "-", -1)
//	transformed = strings.Replace(transformed, " ", "-", -1)
//	return transformed
//}
