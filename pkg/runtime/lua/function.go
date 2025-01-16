package lua

import (
	"github.com/patrixr/glue/pkg/runtime"
	lua "github.com/yuin/gopher-lua"
)

type LuaFunctionVal struct {
	LuaValue[*lua.LFunction]
}

func NewFunc(lv *lua.LFunction) runtime.RTFunction {
	return LuaFunctionVal{
		LuaValue[*lua.LFunction]{lv, runtime.FUNC},
	}
}
