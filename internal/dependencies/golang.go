package dependencies

import (
	"log/slog"
	"strings"

	"github.com/ctroller/goffold/internal/inject"
)

type GoDependencyArgs struct {
	Flags []string
}

func GoDependencyHandler(inject inject.Inject) DependencyHandler {
	return func(dependency Dependency) ([]byte, error) {
		cmdArgs := []string{"get"}
		depName := dependency.Name
		if dependency.Version != "" {
			depName += "@" + dependency.Version
		}

		args, ok := dependency.Args.(GoDependencyArgs)
		if ok {
			flags := args.Flags
			if flags != nil {
				cmdArgs = append(cmdArgs, flags...)
			}
		}

		cmdArgs = append(cmdArgs, depName)

		slog.Debug("Installing dependency...", "dependency", depName)
		out, err := inject.CmdExecutor.Exec("go", cmdArgs...)
		if err != nil {
			return out, err
		}

		slog.Info("Dependency installed successfully.", "dependency", depName)
		slog.Debug("Command executed successfully.", "cmd", "go " + strings.Join(cmdArgs, " "), "output", out)
		
		return out, nil
	}
}