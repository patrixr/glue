package lua

import (
	"github.com/patrixr/glue/pkg/runtime"
	lua "github.com/yuin/gopher-lua"
)

func AnyValue(val lua.LValue) LuaValue[lua.LValue] {
	return LuaValue[lua.LValue]{val, runtime.ANY}
}
