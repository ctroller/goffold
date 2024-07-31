package goffold_test

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

var TestFS = afero.NewMemMapFs()

func OpenMemFile(t *testing.T, name string) afero.File {
	file, err := TestFS.Open(name)
	if err != nil {
		assert.Error(t, err)
	}

	return file
}