package luatools

import lua "github.com/yuin/gopher-lua"

func StrFunc(f func(s string) error) lua.LGFunction {
	return func(L *lua.LState) int {
		L.CheckTypes(1, lua.LTString)

		if err := f(L.ToString(1)); err != nil {
			L.RaiseError(err.Error())
		}

		return 0
	}
}

func Func(f func() error) lua.LGFunction {
	return func(L *lua.LState) int {
		if err := f(); err != nil {
			L.RaiseError(err.Error())
		}

		return 0
	}
}
