package template

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	yaml := `variables:
  - name: name
    type: string
    default: "world"
    prompt: "What is your name?"
  - name: age
    type: int`

	conv, err := ReadVars(strings.NewReader(yaml))
	if err != nil {
		t.Fatal(err)
	}

	assert.Len(t, conv, 2)
	assert.Equal(t, "name", conv[0].Name)
	assert.Equal(t, "string", conv[0].Type)
	assert.Equal(t, "world", conv[0].Default)
	assert.Equal(t, "What is your name?", conv[0].Prompt)
	assert.Equal(t, "age", conv[1].Name)
	assert.Equal(t, "int", conv[1].Type)
	assert.Equal(t, "", conv[1].Default)
	assert.Equal(t, "", conv[1].Prompt)
}
