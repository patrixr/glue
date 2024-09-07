package luatools

import lua "github.com/yuin/gopher-lua"

func GetArgAsString(L *lua.LState, index int) string {
	return L.ToStringMeta(L.Get(1)).String()
}
