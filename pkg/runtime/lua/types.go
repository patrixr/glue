package lua

// import (
// 	"github.com/patrixr/glue/pkg/runtime"
// 	lua "github.com/yuin/gopher-lua"
// )

// type LuaType struct {
// 	name string
// 	lt   lua.LValueType
// }

// func (self LuaType) Is(t runtime.Type) bool {
// 	return ok && lt.lt == self.lt
// }

// func (self LuaType) Name() string {
// 	return self.name
// }

// func (self LuaType) Id() runtime.Type {
// 	return self.id
// }

// var LuaStringType runtime.Type = LuaType{"string", lua.LTString, runtime.STRING}
// var LuaDictType runtime.Type = LuaType{"dict", lua.LTTable, runtime.DICT}
// var LuaNumberType runtime.Type = LuaType{"number", lua.LTNumber, runtime.NUMBER}
// var LuaFunctionType runtime.Type = LuaType{"function", lua.LTFunction, runtime.FUNC}
// var LuaBoolType runtime.Type = LuaType{"bool", lua.LTBool, runtime.BOOL}
// var LuaNilType runtime.Type = LuaType{"nil", lua.LTNil, runtime.NIL}
