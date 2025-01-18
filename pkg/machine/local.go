package machine

import (
	"io"
	"os"
	"os/exec"
	"strings"
)

type LocalMachine struct {
}

func NewLocalMachine() Machine {
	return &LocalMachine{}
}

func (m *LocalMachine) Shell(input string, stdout io.Writer, stderr io.Writer) error {
	args := strings.Fields(input)
	cmd := exec.Command(args[0], args[1:]...)

	cmd.Stdout = stdout
	cmd.Stderr = stderr

	return cmd.Run()
}

func (m *LocalMachine) TempFile(name string) (File, func() error, error) {
	tmp, err := os.CreateTemp("", name)
	return tmp, func() error {
		return os.Remove(tmp.Name())
	}, err
}
