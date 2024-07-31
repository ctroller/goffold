package template

import (
	"testing"

	"github.com/ctroller/goffold/internal/test"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestLoadTemplate(t *testing.T) {
	test.TestFS.MkdirAll("templates/test", 0644)
	afero.WriteFile(test.TestFS, "templates/test/template.yml", []byte(`
template:
  name: test
  description: Test template
  dependencies:
    - name: "go"
      version: "1.16"
      args:
        myarg: "myvalue"
        compiler: [2024]
  layout:
    folders:
      - name: ".devcontainer"
        extends_from: "devcontainers/go"
`), 0644)

	TemplateFs = test.TestFS
	TemplateDir = "templates"
	tpl, err := LoadTemplate("test")
	if err != nil {
		t.Error(err)
		return
	}

	assert.Equal(t, "test", tpl.Name)
	assert.Equal(t, "Test template", tpl.Description)
	assert.Equal(t, "templates/test", tpl.Path)
}
