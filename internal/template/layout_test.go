package template

import (
	"testing"

	"github.com/ctroller/goffold/internal/test"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestLoad(t *testing.T) {
	afero.WriteFile(test.TestFS, "layout.yml", []byte(`
layout:
  folders:
    - name: ".devcontainer"
      extends_from: "devcontainers/go"
    - name: ".devcontainer2"
`), 0644)

	layoutReader := test.OpenMemFile(t, "layout.yml")
	layout, err := LoadLayout(layoutReader)

	if err != nil {
		t.Error(err)
		return
	}

	expected := &Layout{
		Folders: []Folder{
			{
				Name: ".devcontainer",
				ExtendsFrom: "devcontainers/go",
			},
			{
				Name: ".devcontainer2",
			},
		},
	}

	assert.Equal(t, expected, layout)
}