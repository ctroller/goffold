package template

import (
	"fmt"
	"path/filepath"

	"github.com/ctroller/goffold/internal/dependencies"
	"github.com/spf13/afero"
	"gopkg.in/yaml.v3"
)

type Folder struct {
	Name        string `yaml:"name"`
	ExtendsFrom string `yaml:"extends_from"`
}

type Layout struct {
	Folders []Folder `yaml:"folders"`
}

type TemplateVars map[string]string

type Template struct {
	Name         string                    `yaml:"name"`
	Path         string                    `yaml:"-"`
	Description  string                    `yaml:"description"`
	Dependencies []dependencies.Dependency `yaml:"dependencies"`
	Layout       Layout                    `yaml:"layout"`
	Vars         []Var                     `yaml:"variables"`
	TemplateVars TemplateVars              `yaml:"-"`
}

var TemplateDir string
var OutputDir string
var TemplateFs afero.Fs
var OsFs afero.Fs

func ValidateConfig() error {
	_, err := TemplateFs.Stat(TemplateDir)
	if err != nil {
		return fmt.Errorf("templates path %v does not exist (using fs %v)", TemplateDir, TemplateFs)
	}

	return nil
}

func LoadTemplate(name string) (*Template, error) {
	path := filepath.Join(TemplateDir, name)
	_, err := TemplateFs.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("template %v does not exist", name)
	}

	file := filepath.Join(path, "template.yml")
	_, err = TemplateFs.Stat(file)
	if err != nil {
		return nil, fmt.Errorf("template file %v does not exist", file)
	}

	handle, err := TemplateFs.Open(file)
	if err != nil {
		return nil, err
	}
	defer handle.Close()

	var outer struct {
		Template Template `yaml:"template"`
	}
	err = yaml.NewDecoder(handle).Decode(&outer)
	if err != nil {
		return nil, err
	}

	outer.Template.Path = path

	return &outer.Template, nil
}
