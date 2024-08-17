package modules

import (
	"os"

	"github.com/charmbracelet/log"
	"github.com/patrixr/glue/pkg/shell"
	lua "github.com/yuin/gopher-lua"
)

func init() {
	Registry.AddLoader(RegisterLuaPrint)
	Registry.AddLoader(RegisterLuaSh)
}

func RegisterLuaSh(L *lua.LState, _ LuaLifecycle) error {
	L.Register("sh", func(L *lua.LState) int {
		L.CheckTypes(1, lua.LTString)

		cmd := L.ToString(1)
		error := shell.Run(cmd, os.Stdout, os.Stderr)

		if error != nil {
			L.ArgError(1, error.Error())
		}

		return 0
	})

	return nil
}

func RegisterLuaPrint(L *lua.LState, _ LuaLifecycle) error {
	L.SetGlobal("print", L.NewFunction(func(L *lua.LState) int {
		input := L.ToStringMeta(L.Get(1)).String()
		LuaPrint(input)
		return 0
	}))

	return nil
}

func LuaPrint(input string) {
	log.Info(input)
}
