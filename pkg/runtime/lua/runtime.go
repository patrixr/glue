package lua

import (
	"errors"
	"fmt"
	"strings"

	"github.com/patrixr/glue/pkg/runtime"
	lua "github.com/yuin/gopher-lua"
)

type LuaRuntime struct {
	L *lua.LState
}

func NewLuaRuntime() runtime.Runtime {
	L := lua.NewState(lua.Options{
		SkipOpenLibs: true,
	})

	if err := LoadSafeLibs(L); err != nil {
		panic(err.Error())
	}

	return &LuaRuntime{
		L: L,
	}
}

func (luaruntime *LuaRuntime) Close() {
	luaruntime.L.Close()
}

func (luaruntime *LuaRuntime) ExecFile(path string) error {
	return luaruntime.L.DoFile(path)
}

func (luaruntime *LuaRuntime) ExecString(source string) error {
	return luaruntime.L.DoString(source)
}

func (luaruntime *LuaRuntime) String(str string) runtime.RTString {
	return NewString(lua.LString(str))
}

func (luaruntime *LuaRuntime) CheckString(v runtime.RTValue) (runtime.RTString, error) {
	if !v.Type().Is(runtime.STRING) {
		return EmptyString(), fmt.Errorf("Expected a string, received a %s instead", runtime.TypeName(v.Type()))
	}
	str, ok := v.(LuaStringVal)
	if !ok || str.Raw().Type() != lua.LTString {
		return EmptyString(), fmt.Errorf("Expected a string, bad value received")
	}

	return str, nil
}

func (luaruntime *LuaRuntime) CheckNumber(v runtime.RTValue) (runtime.RTNumber, error) {
	if !v.Type().Is(runtime.NUMBER) {
		return EmptyString(), fmt.Errorf("Expected a number, received a %s instead", runtime.TypeName(v.Type()))
	}
	n, ok := v.(LuaNumberVal)
	if !ok || n.Raw().Type() != lua.LTNumber {
		return Zero(), fmt.Errorf("Expected a string, bad value received")
	}
	return n, nil
}

func (luaruntime *LuaRuntime) CheckBool(v runtime.RTValue) (runtime.RTBool, error) {
	if !v.Type().Is(runtime.BOOL) {
		return NewBool(lua.LFalse), fmt.Errorf("Expected a number, received a %s instead", runtime.TypeName(v.Type()))
	}
	n, ok := v.(LuaBoolVal)
	if !ok || n.Raw().Type() != lua.LTBool {
		return NewBool(lua.LFalse), fmt.Errorf("Expected a string, bad value received")
	}
	return n, nil
}

func (luaruntime *LuaRuntime) CheckDict(v runtime.RTValue) (runtime.RTDict, error) {
	if !v.Type().Is(runtime.DICT) {
		return NewDict(nil), fmt.Errorf("Expected a dict, received a %s instead", runtime.TypeName(v.Type()))
	}
	d, ok := v.(LuaDictVal)
	if !ok || d.Raw().Type() != lua.LTTable {
		return NewDict(nil), fmt.Errorf("Expected a dict, bad value received")
	}
	return d, nil
}

func (luaruntime *LuaRuntime) CheckArray(v runtime.RTValue) (runtime.RTArray, error) {
	if !v.Type().Is(runtime.DICT) {
		return NewArray(nil), fmt.Errorf("Expected a dict, received a %s instead", runtime.TypeName(v.Type()))
	}
	d, ok := v.(LuaArrayVal)
	if !ok || d.Raw().Type() != lua.LTTable {
		return NewArray(nil), fmt.Errorf("Expected a dict, bad value received")
	}
	return d, nil
}

func (luaruntime *LuaRuntime) CheckFunction(v runtime.RTValue) (runtime.RTFunction, error) {
	if !v.Type().Is(runtime.FUNC) {
		return NewFunc(nil), fmt.Errorf("Expected a function, received a %s instead", runtime.TypeName(v.Type()))
	}
	fn, ok := v.(LuaFunctionVal)
	if !ok || fn.Raw().Type() != lua.LTFunction {
		return NewFunc(nil), fmt.Errorf("Expected a function, bad value received")
	}
	return fn, nil
}

func (luaruntime *LuaRuntime) EnsureString(v runtime.RTValue) runtime.RTString {
	val, err := luaruntime.CheckString(v)
	if err != nil {
		luaruntime.L.RaiseError("%s", err)
	}
	return val
}

func (luaruntime *LuaRuntime) EnsureNumber(v runtime.RTValue) runtime.RTNumber {
	val, err := luaruntime.CheckNumber(v)
	if err != nil {
		luaruntime.L.RaiseError("%s", err)
	}
	return val
}

func (luaruntime *LuaRuntime) EnsureBool(v runtime.RTValue) runtime.RTBool {
	val, err := luaruntime.CheckBool(v)
	if err != nil {
		luaruntime.L.RaiseError("%s", err)
	}
	return val
}

func (luaruntime *LuaRuntime) EnsureDict(v runtime.RTValue) runtime.RTDict {
	val, err := luaruntime.CheckDict(v)
	if err != nil {
		luaruntime.L.RaiseError("%s", err)
	}
	return val
}

func (luaruntime *LuaRuntime) EnsureFunction(v runtime.RTValue) runtime.RTFunction {
	val, err := luaruntime.CheckFunction(v)
	if err != nil {
		luaruntime.L.RaiseError("%s", err)
	}
	return val
}

func (luaruntime *LuaRuntime) EnsureArray(v runtime.RTValue) runtime.RTArray {
	val, err := luaruntime.CheckArray(v)
	if err != nil {
		luaruntime.L.RaiseError("%s", err)
	}
	return val
}

func (luaruntime *LuaRuntime) RaiseError(format string, args ...interface{}) {
	luaruntime.L.RaiseError(format, args...)
}

func (luaruntime *LuaRuntime) Lang() string {
	return "lua"
}

func (luaruntime *LuaRuntime) SetFunction(
	name string,
	desc string,
	args []runtime.ArgDef,
	fn runtime.CustomFunc,
) error {

	var impl lua.LGFunction = func(L *lua.LState) int {
		var values []runtime.RTValue = []runtime.RTValue{}

		// Based on the expected arguments, we collect all the values form the input
		for i, arg := range args {
			idx := i + 1

			if arg.Type.Is(runtime.STRING) {
				values = append(values, NewString(lua.LString(L.CheckString(idx))))
				continue
			}

			if arg.Type.Is(runtime.NUMBER) {
				values = append(values, NewNumber(L.CheckNumber(idx)))
				continue
			}

			if arg.Type.Is(runtime.DICT) {
				values = append(values, NewDict(L.CheckTable(idx)))
				continue
			}

			if arg.Type.Is(runtime.ARRAY) {
				values = append(values, NewArray(L.CheckTable(idx)))
				continue
			}

			if arg.Type.Is(runtime.FUNC) {
				values = append(values, NewFunc(L.CheckFunction(idx)))
				continue
			}

			if arg.Type.Is(runtime.BOOL) {
				values = append(values, NewBool(lua.LBool(L.CheckBool(idx))))
				continue
			}

			if arg.Type.Is(runtime.ANY) {
				values = append(values, AnyValue(L.CheckAny(idx)))
				continue
			}

			luaruntime.RaiseError("Unsupported type %s for argument %s", runtime.TypeName(arg.Type), arg.Name)
		}

		val := fn(luaruntime, runtime.NewArguments(luaruntime, values))

		if val != nil {
			L.Push(luaruntime.getRawLuaValue(val))
			return 1
		}

		return 0
	}

	var val runtime.RTValue = NewFunc(luaruntime.L.NewFunction(impl))

	_, err := luaruntime.SetGlobal(name, val)

	return err
}

func (luaruntime *LuaRuntime) InvokeFunction(fn runtime.RTFunction, params ...runtime.RTValue) error {
	return luaruntime.invoke(fn, false, params...)
}

func (luaruntime *LuaRuntime) InvokeFunctionSafe(fn runtime.RTFunction, params ...runtime.RTValue) error {
	return luaruntime.invoke(fn, true, params...)
}

func (luaruntime *LuaRuntime) invoke(fn runtime.RTFunction, safe bool, params ...runtime.RTValue) error {
	luafn, ok := fn.(LuaFunctionVal)

	if !ok {
		return fmt.Errorf("Failed to invoke function variable: %s", fn.String())
	}

	var vals []lua.LValue = make([]lua.LValue, len(params))

	// Convert to LValues
	for i, param := range params {
		val, ok := param.(LuaValue[lua.LValue])
		if !ok {
			luaruntime.RaiseError("Invoking function with an invalid argument %s or type %s", param.String(), runtime.TypeName(param.Type()))
		}
		vals[i] = val.Raw()
	}

	return luaruntime.L.CallByParam(lua.P{
		Fn:      luafn.Raw(),
		NRet:    lua.MultRet,
		Protect: safe,
	}, vals...)
}

func (luaruntime *LuaRuntime) SetGlobal(path string, val runtime.RTValue) ([]string, error) {
	L := luaruntime.L
	var value lua.LValue = nil

	switch v := val.(type) {
	case LuaStringVal:
	case LuaNumberVal:
	case LuaDictVal:
	case LuaBoolVal:
	case LuaArrayVal:
	case LuaFunctionVal:
		value = v.Raw()
	default:
		return []string{}, errors.New("SetGlobal failed: Unknown type " + runtime.TypeName(val.Type()))
	}

	keys := strings.Split(path, ".")

	if len(keys) == 0 {
		return []string{}, errors.New("Trying to register a module with empty name")
	}

	if len(keys) == 1 {
		L.SetGlobal(keys[0], value)
		return []string{}, nil
	}

	var ref *lua.LTable
	var err error
	var trail []string

	for i, key := range keys {
		if i == len(keys)-1 {
			break
		}

		trail = append(trail, key)

		if ref == nil {
			// First key, we get the global
			ref, err = luaruntime.GetOrCreateGlobalTable(key)
			if err != nil {
				return []string{}, err
			}
		} else {
			keystr := lua.LString(key)
			nested := ref.RawGet(keystr)

			if nested == lua.LNil {
				table := L.NewTypeMetatable(key)
				ref.RawSet(keystr, table)
				ref = table
				continue
			}

			obj, ok := ref.RawGet(keystr).(*lua.LTable)

			if !ok {
				return []string{}, fmt.Errorf("A table was expected for nested key %s, but instead found %s", key, obj)
			}
			table := L.NewTypeMetatable(key)
			ref.RawSet(keystr, table)
			ref = table
		}
	}

	last := keys[len(keys)-1]
	L.SetField(ref, last, value)

	return trail, nil
}

func (luaruntime *LuaRuntime) GetOrCreateGlobalTable(name string) (*lua.LTable, error) {
	L := luaruntime.L
	ref := L.GetGlobal(name)

	if ref == lua.LNil {
		table := L.NewTypeMetatable(name)
		L.SetGlobal(name, table)
		return table, nil
	}

	table, ok := ref.(*lua.LTable)

	if !ok {
		return nil, fmt.Errorf("Accessing table '%s' failed. Value is not a table", name)
	}

	return table, nil
}

func (luaruntime *LuaRuntime) getRawLuaValue(v runtime.RTValue) lua.LValue {
	switch v.(type) {
	case LuaValue[lua.LValue]:
		return v.(LuaValue[lua.LValue]).Raw()
	case LuaStringVal:
		return lua.LString(v.(LuaStringVal).Raw())
	case LuaBoolVal:
		return lua.LBool(v.(LuaBoolVal).Raw())
	case LuaNumberVal:
		return lua.LNumber(v.(LuaNumberVal).Raw())
	case LuaArrayVal:
		return v.(LuaArrayVal).Raw()
	case LuaDictVal:
		return v.(LuaDictVal).Raw()
	case LuaNilVal:
		return lua.LNil
	default:
		luaruntime.RaiseError("Could not recognize the value as a Lua (%s)", runtime.TypeName(v.Type()))
	}
	return lua.LNil
}
