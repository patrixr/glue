package modules

import (
	"errors"
	"os"
	"path/filepath"
	"time"

	"github.com/patrixr/glue/pkg/core"
	"github.com/patrixr/glue/pkg/luatools"
)

func init() {
	Registry.RegisterModule(func(glue *core.Glue) error {
		glue.Plug().
			Name("backup").
			Short("Creates a backup of a file").
			Long("Creates a backup of a file").
			Do(luatools.StrFunc(func(path string) error {
				return Backup(path)
			}))

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
