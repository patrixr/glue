package lua

import (
	lua "github.com/yuin/gopher-lua"
)

var libs map[string]lua.LGFunction = map[string]lua.LGFunction{
	lua.BaseLibName:   lua.OpenBase,
	lua.TabLibName:    lua.OpenTable,
	lua.StringLibName: lua.OpenString,
	lua.MathLibName:   lua.OpenMath,
}

func LoadSafeLibs(L *lua.LState) error {
	for name, fn := range libs {
		L.Push(L.NewFunction(fn))
		L.Push(lua.LString(name))
		L.Call(1, 0)
	}

	global := L.Get(lua.GlobalsIndex).(*lua.LTable)
	global.RawSetString("collectgarbage", lua.LNil)
	global.RawSetString("dofile", lua.LNil)
	global.RawSetString("load", lua.LNil)
	global.RawSetString("loadfile", lua.LNil)
	global.RawSetString("loadstring", lua.LNil)
	global.RawSetString("setfenv", lua.LNil)
	global.RawSetString("newproxy", lua.LNil)

	return nil
}
