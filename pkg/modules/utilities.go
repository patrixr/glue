package modules

import (
	"os"

	"github.com/patrixr/glue/pkg/core"
	"github.com/patrixr/glue/pkg/luatools"
	"github.com/patrixr/glue/pkg/shell"
	lua "github.com/yuin/gopher-lua"
)

func init() {
	Registry.RegisterModule(UtilitiesMod)
}

func UtilitiesMod(glue *core.Glue) error {
	print := func(L *lua.LState) int {
		input := luatools.GetArgAsString(L, 1)
		glue.Log.Info(input)
		return 0
	}

	sh := luatools.StrFunc(func(cmd string) error {
		return shell.Run(cmd, os.Stdout, os.Stderr)
	})

	glue.Plug().
		Name("sh").
		Short("Run a shell command").
		Long("Run a shell command").
		Example("sh('ls')").
		Do(sh)

	glue.Plug().
		Name("print").
		Short("Print a string").
		Long("Print a string").
		Example("print('Hello, world!')").
		Do(print)

	return nil
}
