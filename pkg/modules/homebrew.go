package modules

import (
	"github.com/patrixr/glue/pkg/homebrew"
	lua "github.com/yuin/gopher-lua"
)

func init() {
	Registry.AddLoader(RegisterHomebrew)
}

func RegisterHomebrew(L *lua.LState, _ LuaLifecycle) error {
	mt := L.NewTypeMetatable("homebrew")
	brew := homebrew.NewHomebrew()

	L.SetGlobal("brew", mt)

	L.SetField(mt, "package", L.NewFunction(func(L *lua.LState) int {
		brew.Brew(L.ToString(1))
		return 0
	}))

	L.SetField(mt, "cask", L.NewFunction(func(L *lua.LState) int {
		brew.Cask(L.ToString(1))
		return 0
	}))

	L.SetField(mt, "tap", L.NewFunction(func(L *lua.LState) int {
		brew.Tap(L.ToString(1))
		return 0
	}))

	L.SetField(mt, "mas", L.NewFunction(func(L *lua.LState) int {
		brew.Mas(L.ToString(1))
		return 0
	}))

	L.SetField(mt, "whalebrew", L.NewFunction(func(L *lua.LState) int {
		brew.Whalebrew(L.ToString(1))
		return 0
	}))

	L.SetField(mt, "sync", L.NewFunction(func(L *lua.LState) int {
		if err := brew.Install(); err != nil {
			L.RaiseError(err.Error())
		}
		return 0
	}))

	return nil
}
