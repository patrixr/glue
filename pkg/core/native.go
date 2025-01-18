package core

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/patrixr/glue/pkg/blueprint"
	. "github.com/patrixr/glue/pkg/runtime"
)

func InstallNativeGlueModules(glue *Glue) {
	glue.Plug("glue.run", FUNCTION).
		Brief("Run a glue script").
		Arg("glue_file", STRING, "the glue file to run").
		Do(func(R Runtime, args *Arguments) (RTValue, error) {
			var resolvedPath string

			file := args.EnsureString(0).String()

			if filepath.IsAbs(file) {
				resolvedPath = file
			} else {
				wd, err := glue.Getwd()
				if err != nil {
					return nil, err
				}
				resolvedPath = filepath.Join(wd, file)
			}

			scriptPath, err := TryFindGlueFile(resolvedPath)

			if err != nil {
				return nil, err
			}

			if err := glue.execFile(scriptPath); err != nil {
				return nil, err
			}
			return nil, nil
		})

	glue.Plug("group", FUNCTION).
		Brief("Create a runnable group").
		Arg("name", STRING, "the name of the group to run").
		Arg("fn", FUNC, "the function to run when the group is invoked").
		Do(func(R Runtime, args *Arguments) (RTValue, error) {
			name := args.EnsureString(0).String()
			fn := args.EnsureFunction(1)

			if len(name) == 0 {
				return nil, errors.New("Group name cannot be empty")
			}

			if name[0] == NegationRune {
				return nil, errors.New(fmt.Sprintf("Group name cannot start with character %c", NegationRune))
			}

			if strings.EqualFold(name, RootLevel) {
				return nil, errors.New(fmt.Sprintf("Group cannot be named %s. Reserved keyword", name))
			}

			if allowed, err := glue.canRunGroup(name); !allowed {
				return nil, err
			}

			glue.Stack.PushGroup(name)

			defer glue.Stack.PopGroup()

			glue.Log.Info("[Group]", "name", name)
			glue.Fire(EV_GROUP_START, name)

			groupPlan := blueprint.NewSerialBlueprint(name)
			basePlan := glue.BluePrint
			glue.BluePrint = groupPlan

			defer func() {
				basePlan.Add(groupPlan)
				glue.BluePrint = basePlan
			}()

			if err := R.InvokeFunctionSafe(fn); err != nil {
				return nil, err
			}

			glue.Fire(EV_GROUP_END, name)

			return nil, nil
		})
}
