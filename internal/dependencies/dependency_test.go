package dependencies

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

var TestFS = afero.NewMemMapFs()

func TestLoad(t *testing.T) {
	afero.WriteFile(TestFS, "dependencies.yml", []byte(`
- name: github.com/ctroller/goffold
  version: latest
  args:
    flags: ["-u"]
- name: github.com/ctroller/goffold2`), 0644)

	afero.WriteFile(TestFS, "dep-versions.yml", []byte(`
versions:
  github.com/ctroller/goffold2: "1.1.1"`), 0644)

	depsReader := openFile(t, "dependencies.yml")
	versionsReader := openFile(t, "dep-versions.yml")

	deps := Load(depsReader, versionsReader)
	expected := []Dependency{
		{
			Name: "github.com/ctroller/goffold",
			Version: "latest",
			Args: map[string]any {
				"flags": []any{"-u"},
			},
		},
		{
			Name: "github.com/ctroller/goffold2",
			Version: "1.1.1",
		},
	}
	
	
	assert.Equal(t, expected, deps)
}

func TestVersionContraints(t *testing.T) {
	afero.WriteFile(TestFS, "dep-versions.yml", []byte(`
versions:
  github.com/ctroller/goffold2: "1.1.1"
  some_package: "latest"`), 0644)

	versionsReader := openFile(t, "dep-versions.yml")
	versions := getVersionConstraints(versionsReader)

	expected := VersionConstraints{
		"github.com/ctroller/goffold2": "1.1.1",
		"some_package": "latest",
	}

	assert.Equal(t, expected, versions)
}

func openFile(t *testing.T, name string) afero.File {
	file, err := TestFS.Open(name)
	if err != nil {
		assert.Error(t, err)
	}

	return file
}