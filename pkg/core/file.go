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

	if file, err := TryFindFile(wd, "glue.lua"); err == nil {
		return file, nil
	}

	// Check for the script in the user's home directory

	glueHome, err := GlueHome()

	if err != nil {
		return "", err
	}

	if file, err := TryFindFile(glueHome, "glue.lua"); err == nil {
		return file, nil
	}

	return "", errors.New("Unable to detect script file")
}

func TryFindGlueFile(path string) (string, error) {
	stats, err := os.Stat(path)

	if os.IsNotExist(err) {
		return "", err
	}

	if stats.IsDir() {
		return TryFindFile(path, "glue.lua")
	}

	return path, nil
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

func GlueHome() (string, error) {
	homedir, err := os.UserHomeDir()

	if err != nil {
		return "", err
	}

	configFolder := filepath.Join(homedir, ".config")
	return filepath.Join(configFolder, "glue"), nil
}
