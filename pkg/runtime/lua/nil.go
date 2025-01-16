package lua

import (
	"github.com/patrixr/glue/pkg/runtime"
	lua "github.com/yuin/gopher-lua"
)

type LuaNilVal struct {
	LuaValue[lua.LValue]
}

func Nil() runtime.RTNil {
	return LuaNilVal{
		LuaValue[lua.LValue]{
			lua.LNil,
			runtime.NIL,
		},
	}
}
