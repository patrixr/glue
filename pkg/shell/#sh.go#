package shell

import (
	"io"
	"os/exec"
	"strings"
)

func Run(input string, stdout io.Writer, stderr io.Writer) error {
	args := strings.Fields(input)
	cmd := exec.Command(args[0], args[1:]...)

	cmd.Stdout = stdout
	cmd.Stderr = stderr

	err := cmd.Start()
	if err != nil {
		return err
	}

	defer cmd.Wait()

	return nil
}
