package template

import (
	"testing"

	"github.com/ctroller/goffold/internal/test"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestLoadTemplate(t *testing.T) {
	test.TestFS.MkdirAll("templates/test", 0644)
	afero.WriteFile(test.TestFS, "templates/test/layout.yml", []byte(`layout:
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
	assert.Equal(t, "templates/test", tpl.Path)
}
