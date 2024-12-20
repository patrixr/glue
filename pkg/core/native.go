package core

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/patrixr/glue/pkg/luatools"
)

const ABOUT_CACHE_KEY = "annotation:about"

func InstallNativeGlueModules(glue *Glue) {
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

	glue.On(EV_NEW_TRACE, func(_ string, data any) error {
		note, ok := glue.Stack.CurrentGroup().Get(ABOUT_CACHE_KEY)

		if !ok || len(note) == 0 {
			return nil
		}

		trace, ok := data.(*Trace)

		if ok {
			trace.About = note
		}

		return nil
	})

	glue.Plug().
		Name("note").
		Short("Annotate the current group with some details").
		Arg("brief", "string", "short explanation of the next step").
		Mode(NONE).
		Bypass().
		Do(luatools.StrFunc(func(s string) error {
			glue.Stack.CurrentGroup().Set(ABOUT_CACHE_KEY, s)
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

			glue.Stack.PushGroup(name)

			defer glue.Stack.PopGroup()

			active, err := glue.AtActiveLevel()

			if err != nil {
				return err
			}

			if active {
				glue.Log.Info("[Group]", "name", name)
				glue.Fire(EV_GROUP_START, name)
			}

			fn()

			glue.Fire(EV_GROUP_END, name)

			return nil
		}))
}
