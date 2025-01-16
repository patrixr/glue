package lua

import (
	"github.com/patrixr/glue/pkg/runtime"
	lua "github.com/yuin/gopher-lua"
)

type LuaValue[T lua.LValue] struct {
	raw T
	typ runtime.Type
}

func (lv LuaValue[T]) Raw() T {
	return lv.raw
}

func (lv LuaValue[T]) Type() runtime.Type {
	return lv.typ
}

func (lv LuaValue[T]) String() string {
	return lv.raw.String()
}
