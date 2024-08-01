package dependencies

import (
	"strings"
	"testing"

	"github.com/ctroller/goffold/internal/inject"
	"github.com/stretchr/testify/assert"
)

var injects = inject.Inject{
	CmdExecutor: inject.CommandExecutor{
		Exec: func(dir string, name string, arg ...string) ([]byte, error) {
			return []byte(strings.Join(append([]string{name}, arg...), " ")), nil
		},
	},
}

func TestSimpleDependency(t *testing.T) {
	dependency := Dependency{
		Pkg: "github.com/ctroller/goffold",
		Args: GoDependencyArgs{
			Flags: []string{"-u"},
		},
	}

	out, err := NewGoResolver(injects).Resolve("", dependency)
	if err != nil {
		assert.Error(t, err)
	}

	assert.Equal(t, "go get -u github.com/ctroller/goffold", string(out))
}
