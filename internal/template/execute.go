package template

import (
	"fmt"
	"log/slog"

	"github.com/spf13/afero"
)

func (t *Template) Execute() error {
	slog.Info("Executing template", "template", t.Name)

	output := OutputDir
	if output == "" {
		output = t.Name
	}

	err := ensureDirDoesNotExist(output)
	if err != nil {
		return err
	}

	err = OsFs.Mkdir(output, 0644)
	if err != nil {
		return err
	}

	return nil
}

func ensureDirDoesNotExist(output string) error {
	exists, _ := afero.DirExists(OsFs, output)
	if exists {
		return fmt.Errorf("output path '%s' already exists", output)
	}

	return nil
}
