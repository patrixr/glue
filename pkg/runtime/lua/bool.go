package lua

import (
	"github.com/patrixr/glue/pkg/runtime"
	lua "github.com/yuin/gopher-lua"
)

type LuaBoolVal struct {
	LuaValue[lua.LBool]
}

func NewBool(lv lua.LBool) runtime.RTBool {
	return LuaBoolVal{
		LuaValue[lua.LBool]{lv, runtime.BOOL},
	}
}

func (lbv LuaBoolVal) Value() bool {
	return bool(lbv.LuaValue.Raw())
}
