package template

import (
	"log/slog"
	"path/filepath"

	"github.com/ctroller/goffold/internal/log"
	"github.com/spf13/afero"
)

type Template struct {
	Name string
	Path string
	Layout Layout
}

var TemplatesPath string
var Fs afero.Fs
var templates map[string]Template

func InitTemplates() error {
	if ok, err := afero.DirExists(Fs, TemplatesPath); !ok {
		log.Fatal("Template path %s does not exist", err, TemplatesPath)
	}

	templates = make(map[string]Template)
	items, err := afero.ReadDir(Fs, TemplatesPath)
	if err != nil {
		log.Fatal("Failed to read templates directory %s", err, TemplatesPath)
	}

	for _, item := range items {
		if item.IsDir() {
			t, err := loadTemplate(item.Name())
			if err != nil {
				slog.Error("Failed to load template", "path", filepath.Join(TemplatesPath, item.Name()), "error", err)
				return err
			} else {
				templates[t.Name] = *t
				slog.Info("Loaded template.", "templateName", t.Name, "path", filepath.Join(TemplatesPath, item.Name()))
			}
		}
	}

	return nil
}

func IsExisting(name string) bool {
	_, ok := templates[name]
	return ok
}

func loadTemplate(name string) (*Template, error) {
	path := filepath.Join(TemplatesPath, name)
	var layout *Layout

	layoutFile := filepath.Join(path, "layout.yml")
	_, err := Fs.Stat(layoutFile)
	if err == nil {
		handle, err := Fs.Open(layoutFile)
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
		Name: name,
		Path: path,
		Layout: l,
	}, nil
}