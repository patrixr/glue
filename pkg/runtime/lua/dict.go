package lua

import (
	"github.com/patrixr/glue/pkg/runtime"
	"github.com/yuin/gluamapper"
	lua "github.com/yuin/gopher-lua"
)

type LuaDictVal struct {
	LuaValue[*lua.LTable]
}

func NewDict(lv *lua.LTable) runtime.RTDict {
	return LuaDictVal{LuaValue[*lua.LTable]{lv, runtime.DICT}}
}

func (dict LuaDictVal) Map() map[interface{}]interface{} {
	opt := gluamapper.Option{
		TagName: "json",
	}
	var tbl lua.LValue = dict.Raw()
	data, ok := gluamapper.ToGoValue(tbl, opt).(map[interface{}]interface{})
	if ok {
		return data
	}
	return map[interface{}]interface{}{}
}
