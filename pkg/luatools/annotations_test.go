package luatools_test

import (
	"strings"
	"testing"

	"github.com/patrixr/glue/pkg/luatools"
)

func TestFunctionAnnotationRender(t *testing.T) {
	funcAnnotation := luatools.LuaFuncAnnotation{
		Name: "myFunction",
		Desc: "This is a test function.\nIt does something useful.",
		Args: []luatools.LuaArgAnnotation{
			{Name: "arg1", Type: "string", Desc: "The first argument."},
			{Name: "arg2", Type: "number", Desc: "The second argument."},
		},
		Returns: []luatools.LuaReturnAnnotation{
			{Type: "boolean", Desc: "Returns true if successful."},
		},
	}

	expected := `---
--- This is a test function.
--- It does something useful.
---
---@param arg1 string The first argument.
---@param arg2 number The second argument.
---
---@return boolean Returns true if successful.
---
function myFunction(arg1, arg2) end

`

	result := funcAnnotation.Render()

	if strings.TrimSpace(result) != strings.TrimSpace(expected) {
		t.Errorf("Render() = %q, want %q", result, expected)
	}
}

func TestLuaClassAnnotation_Render(t *testing.T) {
	classAnnotation := luatools.LuaClassAnnotation{
		Name: "MyClass",
		Fields: []luatools.LuaArgAnnotation{
			{Name: "field1", Type: "string", Desc: "The first field."},
			{Name: "field2", Type: "number", Desc: "The second field."},
		},
	}

	expectedOutput := `---@class MyClass
---@field field1 string The first field.
---@field field2 number The second field.

`

	actualOutput := classAnnotation.Render()

	if actualOutput != expectedOutput {
		t.Errorf("Render() = %v, want %v", actualOutput, expectedOutput)
	}
}
