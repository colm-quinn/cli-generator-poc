package internal

import (
	"log"
	"os"
	"text/template"
)

const (
	commandTemplate = "templates/command.tmpl"
)

//type TemplateActionsData struct {
//	Package   string
//	Store     string
//	SdkModels []string
//	Year      int
//	Command   Command
//	Sdk       Sdk
//	Children  []TemplateActionsData
//}

type Command struct {
	Year           int
	Package        string
	Store          string
	Name           string
	Short          string
	Long           string
	Parameters     []Parameter
	Header         string
	TemplateOutput string
	Children       []Command
}

type Parameter struct {
	Name         string
	Type         string
	IsArray      bool
	Required     bool
	Flag         string
	Usage        string
	Short        string
	DefaultValue string
}

type Sdk struct {
	Model  string
	Store  string
	Method string
}

func Generate(path string, data Command) {
	tmpl, err := template.ParseFiles(commandTemplate)
	if err != nil {
		panic(err)
	}

	path = path + "/" + data.Name

	//TODO: usage, flags, store all generated in their own files and hierarchies
	if len(data.Parameters) == 0 {
		os.MkdirAll(path, 0700)
		//	TODO: Make the root command (different template)
		f, err := os.Create(path + ".go")
		err = tmpl.Execute(f, data)
		if err != nil {
			panic(err)
		}
	} else {
		f, err := os.Create(path + ".go")
		if err != nil {
			log.Println("create file: ", err)
			return
		}

		err = tmpl.Execute(f, data)
		if err != nil {
			panic(err)
		}
	}

	if data.Children != nil {
		for _, child := range data.Children {
			Generate(path, child)
		}
	}
}
