package template

import (
	"io"

	"gopkg.in/yaml.v3"
)

type Folder struct {
	Name        string `yaml:"name"`
	ExtendsFrom string   `yaml:"extends_from"`
}

type Layout struct {
	Folders []Folder `yaml:"folders"`
}

func LoadLayout(in io.Reader) (*Layout, error) {
	var o struct {
		Layout Layout
	}

	err := yaml.NewDecoder(in).Decode(&o)
	if err != nil {
		return nil, err
	}

	return &o.Layout, nil
}