package modules

import (
	"errors"
	"os"
	"path/filepath"
	"time"

	"github.com/patrixr/glue/pkg/core"
	. "github.com/patrixr/glue/pkg/runtime"
)

func init() {
	// @auteur("Modules/Backup")
	//
	// # Backup
	//
	// The `Backup` module provides functionality to create a backup of a specified file. It ensures that the original file is preserved by creating a copy with a timestamp appended to its name.
	//
	// ## Arguments
	//
	// - `path` (string): The path to the file that needs to be backed up.
	//
	// ## Usage
	//
	// To use the `Backup` module, you need to call the `backup` function with the path of the file you want to back up. The function will create a backup file in the same directory with a `.backup.<timestamp>` suffix.
	//
	// ## Example
	//
	// Example usage:
	//
	// ```lua
	// Backup("file.txt")
	// ```
	Registry.RegisterModule(func(glue *core.Glue) error {
		glue.Plug("backup", core.MODULE).
			Brief("Creates a backup of a file").
			Arg("path", STRING, "the file to create a backup of").
			Do(func(R Runtime, args *Arguments) (RTValue, error) {
				str := args.EnsureString(0)
				return nil, Backup(str.String())
			})

		return nil
	})
}

func Backup(originalPath string) error {
	stat, err := os.Stat(originalPath)

	if err != nil {
		return err
	}

	if stat.IsDir() {
		return errors.New("Cannot create backup of a folder: " + originalPath)
	}

	dir := filepath.Dir(originalPath)
	name := filepath.Base(originalPath)
	backupName := name + ".backup." + time.Now().Format(time.RFC3339)
	backupPath := filepath.Join(dir, backupName)

	content, err := os.ReadFile(originalPath)

	if err != nil {
		return err
	}

	if err = os.WriteFile(backupPath, content, stat.Mode()); err != nil {
		return err
	}

	return nil
}
