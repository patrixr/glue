package machine

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"time"
)

type HomebrewParams struct {
	Packages   []string `json:"packages"`
	Casks      []string `json:"casks"`
	Taps       []string `json:"taps"`
	Mas        []string `json:"mas"`
	Whalebrews []string `json:"whalebrews"`
}

type Row struct {
	kind string
	name string
}

type Homebrew struct {
	machine Machine
	rows    []Row
	stdout  io.Writer
	stderr  io.Writer
}

func HomebrewBundle(m Machine, params HomebrewParams, stdout io.Writer, stderr io.Writer) error {
	tmp, close, err := m.TempFile(".glue_brewfile_" + time.Now().Format("20060102150405"))

	if err != nil {
		return err
	}

	defer close()

	for _, row := range params.Packages {
		tmp.Write([]byte(fmt.Sprintf("brew \"%s\"\n", row)))
	}

	for _, row := range params.Casks {
		tmp.Write([]byte(fmt.Sprintf("cask \"%s\"\n", row)))
	}

	for _, row := range params.Taps {
		tmp.Write([]byte(fmt.Sprintf("tap \"%s\"\n", row)))
	}

	for _, row := range params.Mas {
		tmp.Write([]byte(fmt.Sprintf("mas \"%s\"\n", row)))
	}

	for _, row := range params.Whalebrews {
		tmp.Write([]byte(fmt.Sprintf("whalebrew \"%s\"\n", row)))
	}

	return m.Shell(fmt.Sprintf("brew bundle --file=%s --no-lock", tmp.Name()), stdout, stderr)
}

func HomebrewUpgrade(m Machine, stdout io.Writer, stderr io.Writer) error {
	return m.Shell("brew upgrade", stdout, stderr)
}

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

func IsHomebrewInstalled(m Machine) bool {
	path, err := GetHomebrewBin()
	return err == nil && path != ""
}

func InstallHomebrew(m Machine, stdout io.Writer, stderr io.Writer) error {
	if runtime.GOOS != "darwin" && runtime.GOOS != "linux" {
		return errors.New("Homebrew is only supported on macOS and Linux")
	}

	if IsHomebrewInstalled(m) {
		stdout.Write([]byte("Homebrew detected on system\n"))
		return nil
	}

	installCommand := `/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"`

	return m.Shell(
		fmt.Sprintf("bash -c \"%s\"", installCommand),
		stdout,
		stderr,
	)
}

func UpdateHomebrew(m Machine, stdout io.Writer, stderr io.Writer) error {
	path, err := GetHomebrewBin()

	if err != nil {
		return err
	}

	return m.Shell(path+" update", stdout, stderr)
}
