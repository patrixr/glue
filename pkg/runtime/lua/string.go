package lua

import (
	"github.com/patrixr/glue/pkg/runtime"
	lua "github.com/yuin/gopher-lua"
)

type LuaStringVal struct {
	LuaValue[lua.LString]
}

func NewString(lv lua.LString) runtime.RTString {
	return LuaStringVal{
		LuaValue[lua.LString]{lv, runtime.STRING},
	}
}

func EmptyString() runtime.RTString {
	return NewString(lua.LString(""))
}
