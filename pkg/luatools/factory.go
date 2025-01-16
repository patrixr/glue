package luatools

import (
	"github.com/yuin/gluamapper"
	lua "github.com/yuin/gopher-lua"
)

type Callback func() error

type RaiseErrorFunc func(string, ...interface{})

type LuaFuncWithError func(*lua.LState) (int, error)

func StrFunc(f func(s string) error) LuaFuncWithError {
	return func(L *lua.LState) (int, error) {
		L.CheckTypes(1, lua.LTString)

		if err := f(L.ToString(1)); err != nil {
			return 0, err
		}

		return 0, nil
	}
}

func Str2Func(f func(s1 string, s2 string) error) LuaFuncWithError {
	return func(L *lua.LState) (int, error) {
		L.CheckTypes(1, lua.LTString)
		L.CheckTypes(2, lua.LTString)

		if err := f(L.ToString(1), L.ToString(2)); err != nil {
			return 0, err
		}

		return 0, nil
	}
}

func StrInStrOutFunc(f func(s string) (string, error)) LuaFuncWithError {
	return func(L *lua.LState) (int, error) {
		L.CheckTypes(1, lua.LTString)

		out, err := f(L.ToString(1))

		if err != nil {
			L.Push(lua.LString(""))
			return 1, err
		}

		L.Push(lua.LString(out))

		return 1, nil
	}
}

func StrFuncWithOpts[T any](f func(s string, opts T) error) LuaFuncWithError {
	return func(L *lua.LState) (int, error) {
		L.CheckTypes(1, lua.LTString)
		L.CheckTypes(2, lua.LTTable)

		var data T

		if err := gluamapper.Map(L.ToTable(2), &data); err != nil {
			return 0, err
		}

		if err := f(L.ToString(1), data); err != nil {
			return 0, err
		}

		return 0, nil
	}
}

func TableFunc[T any](f func(params T) error) LuaFuncWithError {
	return func(L *lua.LState) (int, error) {
		L.CheckTypes(1, lua.LTTable)

		var data T

		if err := gluamapper.Map(L.ToTable(1), &data); err != nil {
			return 0, err
		}

		if err := f(data); err != nil {
			return 0, err
		}

		return 0, nil
	}
}

func Func(f func() error) LuaFuncWithError {
	return func(L *lua.LState) (int, error) {
		if err := f(); err != nil {
			return 0, err
		}

		return 0, nil
	}
}

func CallbackFunc(f func(cb func() error) error) LuaFuncWithError {
	return func(L *lua.LState) (int, error) {
		L.CheckTypes(1, lua.LTFunction)

		invoke := func() error {
			return L.CallByParam(lua.P{
				Fn:      L.ToFunction(1),
				NRet:    1,
				Protect: true,
			})
		}

		if err := f(invoke); err != nil {
			return 0, err
		}

		return 0, nil
	}
}

func NamedCallbackFunc(f func(name string, cb Callback) error) LuaFuncWithError {
	return func(L *lua.LState) (int, error) {
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
			return 0, err
		}

		return 0, nil
	}
}

func MockFunc(args ...lua.LValue) LuaFuncWithError {
	return func(L *lua.LState) (int, error) {
		for _, arg := range args {
			L.Push(arg)
		}

		return len(args), nil
	}
}

func UnimplementedFunc(errMessage string) LuaFuncWithError {
	return func(L *lua.LState) (int, error) {
		L.RaiseError("%s", errMessage)
		return 0, nil
	}
}
