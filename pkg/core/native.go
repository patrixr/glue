package core

import (
	"path/filepath"

	"github.com/patrixr/glue/pkg/luatools"
)

func installNativeGlueModules(glue *Glue) {

	glue.Plug().
		Name("glue.run").
		Short("Run a glue script").
		Do(luatools.StrFunc(func(file string) error {
			var resolvedPath string

			if filepath.IsAbs(file) {
				resolvedPath = file
			} else {
				wd, err := glue.Getwd()
				if err != nil {
					return err
				}
				resolvedPath = filepath.Join(wd, file)
			}

			scriptPath, err := TryFindGlueFile(resolvedPath)

			if err != nil {
				return err
			}

			if err := glue.RunFileRaw(scriptPath); err != nil {
				return err
			}
			return nil
		}))
}
