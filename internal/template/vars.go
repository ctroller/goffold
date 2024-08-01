package template

import (
	"bytes"
	"io"
	"text/template"

	"github.com/manifoldco/promptui"
	"gopkg.in/yaml.v3"
)

type Var struct {
	Name    string
	Type    string
	Default string
	Prompt  string
}

func (v Var) GetValue(t *Template) (string, error) {
	tpl := template.Must(template.New("").Parse(v.Default))
	var buf bytes.Buffer

	err := tpl.Execute(&buf, t)
	if err != nil {
		return "", err
	}
	def := buf.String()

	prompt := promptui.Prompt{
		Label:   v.Prompt,
		Default: def,
	}

	return prompt.Run()
}

func ReadVars(in io.Reader) ([]Var, error) {
	var outer struct {
		Variables []Var
	}

	err := yaml.NewDecoder(in).Decode(&outer)
	return outer.Variables, err
}
