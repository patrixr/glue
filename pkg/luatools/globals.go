package luatools

import (
	"errors"
	"fmt"
	"strings"

	lua "github.com/yuin/gopher-lua"
)

func GetOrCreateGlobalTable(L *lua.LState, name string) (*lua.LTable, error) {
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

// Sets a value inside nested tables in the global scope
// Will create nested tables as needed to create a tree structure
// as specified by the path argument
func SetNestedGlobalValue(L *lua.LState, path string, value lua.LValue) ([]string, error) {
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
			ref, err = GetOrCreateGlobalTable(L, key)
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
