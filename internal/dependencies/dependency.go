package dependencies

import (
	"io"
	"os"

	"github.com/ctroller/goffold/internal/log"
	"gopkg.in/yaml.v3"
)

type Dependency struct {
	Name    string
	Version string `default:"latest"`
	Args    any
}

type VersionConstraints map[string]string

func Load(depsReader io.Reader, versionsReader io.Reader) []Dependency {
	versions := getVersionConstraints(versionsReader)

	decoder := yaml.NewDecoder(depsReader)
	var deps []Dependency
	err := decoder.Decode(&deps)
	if err != nil {
		log.Fatal("Failed to decode dependencies", err)
		os.Exit(1)
	}

	if versions != nil {
		for i, dep := range deps {
			if v, ok := versions[dep.Name]; ok {
				deps[i].Version = v
			}
		}
	}

	return deps
}

func getVersionConstraints(in io.Reader) VersionConstraints {
	if in == nil {
		return nil
	}

	var outer struct {
		Versions VersionConstraints `yaml:"versions"`
	}

	decoder := yaml.NewDecoder(in)
	err := decoder.Decode(&outer)
	if err != nil {
		log.Fatal("Failed to decode dependencies", err)
	}

	return outer.Versions
}