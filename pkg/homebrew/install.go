package homebrew

import (
	"errors"
	"fmt"

	"os/exec"
	"runtime"

	"io"
	"os"

	"github.com/patrixr/glue/pkg/shell"
)

func GetHomebrewBin() (string, error) {
	cmd := exec.Command("which", "brew")
	out, err := cmd.Output()

	if err != nil {
		return string(out), nil
	}

	paths := []string{
		"/opt/homebrew/bin/brew",              // Apple Silicon Macs
		"/usr/local/bin/brew",                 // Common path on macOS
		"/home/linuxbrew/.linuxbrew/bin/brew", // Common path on Linux
	}

	for _, path := range paths {
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}

	return "", errors.New("Homebrew not found")
}

func IsHomebrewInstalled() bool {
	path, err := GetHomebrewBin()
	return err == nil && path != ""
}

func InstallHomebrew(stdout io.Writer, stderr io.Writer) error {
	if runtime.GOOS != "darwin" && runtime.GOOS != "linux" {
		return errors.New("Homebrew is only supported on macOS and Linux")
	}

	if IsHomebrewInstalled() {
		stdout.Write([]byte("Homebrew detected on system\n"))
		return nil
	}

	installCommand := `/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"`

	return shell.Run(
		fmt.Sprintf("bash -c \"%s\"", installCommand),
		stdout,
		stderr,
	)
}

func UpdateHomebrew(stdout io.Writer, stderr io.Writer) error {
	path, err := GetHomebrewBin()

	if err != nil {
		return err
	}

	return shell.Run(path+" update", stdout, stderr)
}
