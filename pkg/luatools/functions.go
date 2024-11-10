package luatools

import (
	"github.com/yuin/gluamapper"
	lua "github.com/yuin/gopher-lua"
)

type Callback func() error

func StrFunc(f func(s string) error) lua.LGFunction {
	return func(L *lua.LState) int {
		L.CheckTypes(1, lua.LTString)

		if err := f(L.ToString(1)); err != nil {
			L.RaiseError(err.Error())
		}

		return 0
	}
}

func Str2Func(f func(s1 string, s2 string) error) lua.LGFunction {
	return func(L *lua.LState) int {
		L.CheckTypes(1, lua.LTString)
		L.CheckTypes(2, lua.LTString)

		if err := f(L.ToString(1), L.ToString(2)); err != nil {
			L.RaiseError(err.Error())
		}

		return 0
	}
}

func StrInStrOutFunc(f func(s string) (string, error)) lua.LGFunction {
	return func(L *lua.LState) int {
		L.CheckTypes(1, lua.LTString)

		out, err := f(L.ToString(1))

		if err != nil {
			L.RaiseError(err.Error())
		}

		L.Push(lua.LString(out))

		return 1
	}
}

func StrFuncWithOpts[T any](f func(s string, opts T) error) lua.LGFunction {
	return func(L *lua.LState) int {
		L.CheckTypes(1, lua.LTString)
		L.CheckTypes(2, lua.LTTable)

		var data T

		if err := gluamapper.Map(L.ToTable(2), &data); err != nil {
			L.RaiseError(err.Error())
		}

		if err := f(L.ToString(1), data); err != nil {
			L.RaiseError(err.Error())
		}

		return 0
	}
}

func TableFunc[T any](f func(params T) error) lua.LGFunction {
	return func(L *lua.LState) int {
		L.CheckTypes(1, lua.LTTable)

		var data T

		if err := gluamapper.Map(L.ToTable(1), &data); err != nil {
			L.RaiseError(err.Error())
		}

		if err := f(data); err != nil {
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

func CallbackFunc(f func(cb func() error) error) lua.LGFunction {
	return func(L *lua.LState) int {
		L.CheckTypes(1, lua.LTFunction)

		invoke := func() error {
			return L.CallByParam(lua.P{
				Fn:      L.ToFunction(1),
				NRet:    1,
				Protect: true,
			})
		}

		if err := f(invoke); err != nil {
			L.RaiseError(err.Error())
		}

		return 0
	}
}

func NamedCallbackFunc(f func(name string, cb Callback) error) lua.LGFunction {
	return func(L *lua.LState) int {
		L.CheckTypes(1, lua.LTString)
		L.CheckTypes(2, lua.LTFunction)

		invoke := func() error {
			return L.CallByParam(lua.P{
				Fn:      L.ToFunction(2),
				NRet:    1,
				Protect: true,
			})
		}

		if err := f(L.ToString(1), invoke); err != nil {
			L.RaiseError(err.Error())
		}

		return 0
	}
}
