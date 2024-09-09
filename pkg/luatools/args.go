package luatools

import lua "github.com/yuin/gopher-lua"

func GetArgAsString(L *lua.LState, index int) string {
	return L.ToStringMeta(L.Get(index)).String()
}

func GetAllArgsAsStrings(L *lua.LState) []string {
	inputs := []string{}
	idx := 1

	for {
		arg := L.Get(idx)

		if arg == lua.LNil {
			break
		}

		inputs = append(inputs, L.ToStringMeta(arg).String())
		idx++
	}
	return inputs
}
