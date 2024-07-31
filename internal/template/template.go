package template

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/afero"
)

type Structure struct {
}

type Template struct {
	Name   string
	Path   string
	Layout Layout
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
	var layout *Layout

	layoutFile := filepath.Join(path, "layout.yml")
	_, err := TemplateFs.Stat(layoutFile)
	if err == nil {
		handle, err := TemplateFs.Open(layoutFile)
		if err != nil {
			return nil, err
		}
		defer handle.Close()

		layout, err = LoadLayout(handle)
		if err != nil {
			return nil, err
		}
	}

	var l Layout
	if layout != nil {
		l = *layout
	}

	return &Template{
		Name:   name,
		Path:   path,
		Layout: l,
	}, nil
}
