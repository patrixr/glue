package core

import (
	"errors"
	"os"
	"path/filepath"
)

func AutoDetectScriptFile() (string, error) {
	// Check for the script in the current working directory

	wd, err := os.Getwd()

	if err != nil {
		return "", err
	}

	if file, err := TryFindFile(wd, ".glue.lua"); err == nil {
		return file, nil
	}

	if file, err := TryFindFile(wd, "glue.lua"); err == nil {
		return file, nil
	}

	// Check for the script in the user's home directory

	homeDir, err := os.UserHomeDir()

	if err != nil {
		return "", err
	}

	if file, err := TryFindFile(homeDir, ".glue.lua"); err == nil {
		return file, nil
	}

	if file, err := TryFindFile(homeDir, "glue.lua"); err == nil {
		return file, nil
	}

	// Check for the script in the user's configuration directory

	configFolder := filepath.Join(homeDir, ".config", "glue")

	if file, err := TryFindFile(configFolder, "glue.lua"); err == nil {
		return file, nil
	}

	return "", errors.New("Unable to detect script file")
}

func TryFindFile(directory string, filename string) (string, error) {
	_, err := os.Stat(directory)

	if os.IsNotExist(err) {
		return "", err
	}

	file := filepath.Join(directory, filename)
	_, err = os.Stat(file)

	if err == nil {
		return file, nil
	}

	return "", err
}
