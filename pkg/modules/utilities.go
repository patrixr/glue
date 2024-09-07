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
	sh := luatools.StrFunc(func(cmd string) error {
		return shell.Run(cmd, os.Stdout, os.Stderr)
	})

	print := func(L *lua.LState) int {
		input := luatools.GetArgAsString(L, 1)
		glue.Log.Info(input)
		return 0
	}

	glue.AddFunction("sh", sh)
	glue.AddFunction("print", print)

	return nil
}
