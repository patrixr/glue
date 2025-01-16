package lua

import (
	"github.com/patrixr/glue/pkg/runtime"
	lua "github.com/yuin/gopher-lua"
)

type LuaNumberVal struct {
	LuaValue[lua.LNumber]
}

func NewNumber(lv lua.LNumber) runtime.RTNumber {
	return LuaNumberVal{
		LuaValue[lua.LNumber]{lv, runtime.NUMBER},
	}
}

func Zero() runtime.RTNumber {
	var zero lua.LNumber = 0
	return NewNumber(zero)
}
