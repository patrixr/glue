package lua

import (
	"github.com/patrixr/glue/pkg/runtime"
	"github.com/yuin/gluamapper"
	lua "github.com/yuin/gopher-lua"
)

type LuaArrayVal struct {
	LuaValue[*lua.LTable]
}

func NewArray(lv *lua.LTable) runtime.RTArray {
	return LuaArrayVal{LuaValue[*lua.LTable]{lv, runtime.ARRAY}}
}

func (dict LuaArrayVal) Map() []interface{} {
	opt := gluamapper.Option{
		TagName: "json",
	}
	var tbl lua.LValue = dict.Raw()
	data, ok := gluamapper.ToGoValue(tbl, opt).([]interface{})
	if ok {
		return data
	}
	return data
}
