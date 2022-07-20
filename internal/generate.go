package internal

import (
	"os"
	"text/template"
)

const (
	generatedPath   = "generated"
	commandTemplate = "templates/command.tmpl"
)

type TemplateActionsData struct {
	Package   string
	Store     string
	SdkModels []string
	Year      int
	Command   Command
}

type Command struct {
	Name           string
	Short          string
	Long           string
	Parameters     []Parameter
	Header         string
	TemplateOutput string
}

type Parameter struct {
	Name     string
	Type     string
	Required bool
}

func Generate(data TemplateActionsData) {
	tmpl, err := template.ParseFiles(commandTemplate)
	if err != nil {
		panic(err)
	}

	err = tmpl.Execute(os.Stdout, data)
	if err != nil {
		panic(err)
	}

	// TODO: Write to file
}
