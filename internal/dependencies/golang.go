package dependencies

import (
	"log/slog"
	"strings"

	"github.com/ctroller/goffold/internal/inject"
)

type GoDependencyArgs struct {
	Flags []string
}

func NewGoResolver(inject inject.Inject) DependencyResolver {
	return DependencyResolver{
		Type: "go",
		Resolve: func(dir string, dependency Dependency) ([]byte, error) {
			cmdArgs := []string{"get"}
			depName := dependency.Pkg
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

			slog.Info("Installing dependency...", "dependency", depName)
			out, err := inject.CmdExecutor.Exec(dir, "go", cmdArgs...)
			if err != nil {
				return out, err
			}

			slog.Info("Dependency installed successfully.", "dependency", depName)
			slog.Debug("Command executed successfully.", "cmd", "go "+strings.Join(cmdArgs, " "), "output", out)

			return out, nil
		},
		Finisher: func(dir string) error {
			slog.Info("Running go mod tidy...")
			_, err := inject.CmdExecutor.Exec(dir, "go", "mod", "tidy")
			return err
		},
	}
}
