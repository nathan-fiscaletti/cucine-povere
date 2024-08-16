package generator

import (
	"fmt"
	"html/template"
	"os"
)

// FillNamedTemplate fills the template with the data and writes it to the output file
// with the given name
func FillNamedTemplate(input string, name string, data interface{}) error {
	out, err := os.Create(fmt.Sprintf("./public/%s", name))
	if err != nil {
		return err
	}

	tmpl, err := template.ParseFiles(fmt.Sprintf("./templates/%s", input))
	if err != nil {
		return err
	}

	err = tmpl.Execute(out, data)
	if err != nil {
		return err
	}

	return nil
}

// FillTemplate fills the template with the data and writes it to the output file
func FillTemplate(name string, data interface{}) error {
	return FillNamedTemplate(name, name, data)
}
