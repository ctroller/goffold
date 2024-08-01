package inject

import (
	"os/exec"
)

type CommandExecutor struct {
	Exec func(dir string, name string, arg ...string) ([]byte, error)
}

type Inject struct {
	CmdExecutor CommandExecutor
}

var Defaults = Inject{
	CmdExecutor: DefaultCommandExecutor(),
}

func DefaultCommandExecutor() CommandExecutor {
	return CommandExecutor{
		Exec: func(dir string, name string, arg ...string) ([]byte, error) {
			cmdExec := exec.Command(name, arg...)
			cmdExec.Dir = dir
			return cmdExec.Output()
		},
	}
}
