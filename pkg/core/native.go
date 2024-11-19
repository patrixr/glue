package core

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/patrixr/glue/pkg/luatools"
)

func InstallNativeGlueModules(glue *Glue) {
	pop := func(s []string) []string {
		if len(s) == 0 {
			return s
		}
		return s[:len(s)-1]
	}

	glue.Plug().
		Name("glue").
		Short("Run a glue script").
		Arg("glue_file", "string", "the glue file to run").
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

	glue.Plug().
		Name("group").
		Short("Create a runnable group").
		Arg("name", "string", "the name of the group to run").
		Arg("fn", "function", "the function to run when the group is invoked").
		Mode(NONE).
		Bypass().
		Do(luatools.NamedCallbackFunc(func(name string, fn luatools.Callback) error {
			if len(name) == 0 {
				return errors.New("Group name cannot be empty")
			}

			if name[0] == NegationRune {
				return errors.New(fmt.Sprintf("Group name cannot start with character %c", NegationRune))
			}

			if strings.EqualFold(name, RootLevel) {
				return errors.New(fmt.Sprintf("Group cannot be named %s. Reserved keyword", name))
			}

			glue.nesting = append(glue.nesting, name)

			defer func() {
				glue.nesting = pop(glue.nesting)
			}()

			active, err := glue.AtActiveLevel()

			if err != nil {
				return err
			}

			if active {
				glue.Log.Info("[Group]", "path", strings.Join(glue.nesting[1:], "."))
			}

			fn()

			return nil
		}))
}
