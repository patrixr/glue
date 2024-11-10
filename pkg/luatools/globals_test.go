package luatools

import (
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	lua "github.com/yuin/gopher-lua"
)

func TestSetNestedGlobalValue(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	tests := []struct {
		name  string
		path  string
		value lua.LValue
		err   error
	}{
		{
			name:  "SimpleGlobal",
			path:  "myvar",
			value: lua.LString("hello"),
			err:   nil,
		},
		{
			name:  "NestedTable",
			path:  "mytable.nested.var",
			value: lua.LNumber(42),
			err:   nil,
		},
		{
			name:  "EmptyPath",
			path:  "",
			value: lua.LString("hello"),
			err:   errors.New("Trying to register a module with empty name"),
		},
		{
			name:  "DeeplyNested",
			path:  "a.b.c.d.e.f",
			value: lua.LBool(true),
			err:   nil,
		},
		{
			name:  "NilValue",
			path:  "nil.test",
			value: lua.LNil,
			err:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := SetNestedGlobalValue(L, tt.path, tt.value)
			assert.NoError(t, err)

			// Verify that the value was set correctly
			keys := strings.Split(tt.path, ".")
			ref := L.GetGlobal(keys[0])
			for i := 1; i < len(keys); i++ {
				if ref == lua.LNil {
					t.Errorf("Expected a table at path '%s', but found nil", strings.Join(keys[:i+1], "."))
					return
				}

				table, ok := ref.(*lua.LTable)
				if !ok {
					t.Errorf("Expected a table at path '%s', but found %T", strings.Join(keys[:i+1], "."), ref)
					return
				}
				ref = table.RawGet(lua.LString(keys[i]))
			}

			assert.Equal(t, ref, tt.value)

		})
		L.SetTop(0) // Clear the stack after each test case. Important for "ExistingNonTable" test case
	}

	t.Run("Nesting a value on a non-table should fail", func(t *testing.T) {
		L := lua.NewState()
		defer L.Close()

		_, err := SetNestedGlobalValue(L, "rootkey", lua.LString("hi"))
		assert.NoError(t, err)
		_, err = SetNestedGlobalValue(L, "rootkey.subkey", lua.LString("hi again"))
		assert.Error(t, err)
	})
}
