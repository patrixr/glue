package core

import (
	"path/filepath"

	"github.com/patrixr/glue/pkg/luatools"
)

func installNativeGlueModules(glue *Glue) {
	glue.AddFunction("glue.run", luatools.StrFunc(func(file string) error {
		var resolvedPath string

		if filepath.IsAbs(file) {
			resolvedPath = file
		} else {
			current, err := glue.GetCurrentScript()
			if err != nil {
				return err
			}
			resolvedPath = filepath.Join(filepath.Dir(current), file)
		}

		scriptPath, err := TryFindGlueFile(resolvedPath)

		if err != nil {
			return err
		}

		if err := glue.RunRaw(scriptPath); err != nil {
			return err
		}
		return nil
	}))
}
