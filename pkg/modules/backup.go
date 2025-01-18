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
