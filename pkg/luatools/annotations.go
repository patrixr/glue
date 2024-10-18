package luatools

import (
	"fmt"
	"strings"

	"github.com/patrixr/q"
)

// Function annotations

type LuaArgAnnotation struct {
	Name string
	Type string
	Desc string
}

type LuaReturnAnnotation struct {
	Type string
	Desc string
}

type LuaFuncAnnotation struct {
	Name    string
	Desc    string
	Args    []LuaArgAnnotation
	Returns []LuaReturnAnnotation
}

func (funcAnnotation *LuaFuncAnnotation) Render() string {
	var builder strings.Builder

	builder.WriteString("---\n")

	for _, line := range strings.Split(funcAnnotation.Desc, "\n") {
		builder.WriteString(fmt.Sprintf("--- %s\n", line))
	}

	builder.WriteString("---\n")

	for _, arg := range funcAnnotation.Args {
		builder.WriteString(fmt.Sprintf("---@param %s %s %s\n", arg.Name, arg.Type, arg.Desc))
	}

	builder.WriteString("---\n")

	for _, ret := range funcAnnotation.Returns {
		builder.WriteString(fmt.Sprintf("---@return %s %s\n", ret.Type, ret.Desc))
	}

	builder.WriteString("---\n")

	builder.WriteString(
		fmt.Sprintf("function %s(%s) end",
			funcAnnotation.Name,
			strings.Join(q.Map(funcAnnotation.Args,
				func(arg LuaArgAnnotation) string {
					return arg.Name
				}), ", "),
		),
	)

	builder.WriteString("\n")

	return builder.String()
}

// Class annotations

type LuaClassAnnotation struct {
	Name   string
	Fields []LuaArgAnnotation
}

func (classAnnotation *LuaClassAnnotation) Render() string {
	var builder strings.Builder

	builder.WriteString(fmt.Sprintf("---@class %s\n", classAnnotation.Name))

	for _, field := range classAnnotation.Fields {
		builder.WriteString(fmt.Sprintf("---@field %s %s %s\n", field.Name, field.Type, field.Desc))
	}

	builder.WriteString("\n")

	return builder.String()
}
