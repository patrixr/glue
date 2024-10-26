package luatools

import (
	"fmt"
	"strings"

	"github.com/patrixr/q"
)

const (
	META  = "meta"
	FUNC  = "function"
	CLASS = "class"
	TABLE = "table"
)

type Annotable interface {
	Render() string
	Type() string
}

//
// Annotation collection
//

type LuaAnnotations struct {
	items []Annotable
}

func (annotations *LuaAnnotations) Render() string {
	var builder strings.Builder

	builder.WriteString("--@meta\n")

	for _, annotation := range annotations.items {
		builder.WriteString(annotation.Render())
		builder.WriteRune('\n')
	}

	return builder.String()
}

func (annotations *LuaAnnotations) Add(annotation Annotable) {
	annotations.items = append(annotations.items, annotation)
}

func (annotations *LuaAnnotations) AddClass(name string) *LuaClassAnnotation {
	class := &LuaClassAnnotation{
		Name: name,
	}
	annotations.items = append(annotations.items, class)
	return class
}

func (annotations *LuaAnnotations) Type() string {
	return META
}

func (annotations *LuaAnnotations) FindAllByType(kind string) []Annotable {
	return q.FindAll(annotations.items, func(a Annotable, _ int) bool {
		return a.Type() == kind
	})
}

func (annotations *LuaAnnotations) FindTable(name string) *LuaTableAnnotation {
	tables := annotations.FindAllByType(TABLE)
	for _, item := range tables {
		table, ok := item.(*LuaTableAnnotation)

		if !ok {
			continue
		}

		if table.Name == name {
			return table
		}
	}
	return nil
}

func (annotations *LuaAnnotations) AddNestedTable(path []string) {
	if len(path) == 0 {
		return
	}

	head := path[0]
	tail := path[1:]
	base := annotations.FindTable(head)

	if base == nil {
		base = &LuaTableAnnotation{
			Name: head,
		}
		annotations.Add(base)
	}

	for _, key := range tail {
		found, child, _ := q.Find(base.Children, func(t *LuaTableAnnotation, _ int) bool {
			return t.Name == key
		})

		if found {
			base = child
			continue
		}

		base = base.AddChild(key)
	}
}

// (Nested) Table annotations
type LuaTableAnnotation struct {
	Name     string
	Children []*LuaTableAnnotation
}

func (tableAnno *LuaTableAnnotation) Type() string {
	return TABLE
}

func (tableAnno *LuaTableAnnotation) Render() string {
	var traverse func(ref *LuaTableAnnotation, level int)
	var builder strings.Builder

	traverse = func(ref *LuaTableAnnotation, level int) {
		builder.WriteString(fmt.Sprintf("%s%s = {\n", strings.Repeat("  ", level), ref.Name))
		for _, child := range ref.Children {
			traverse(child, level+1)
		}
		builder.WriteString(fmt.Sprintf("%s}\n", strings.Repeat("  ", level)))
	}

	traverse(tableAnno, 0)

	return builder.String()
}

func (tableAnno *LuaTableAnnotation) AddChild(name string) *LuaTableAnnotation {
	child := &LuaTableAnnotation{
		Name: name,
	}
	tableAnno.Children = append(tableAnno.Children, child)
	return child
}

//
// Function annotations
//

type LuaFieldDesc struct {
	Name string
	Type string
	Desc string
}

type LuaReturnDesc struct {
	Type string
	Desc string
}

type LuaFuncAnnotation struct {
	Name    string
	Desc    string
	Args    []LuaFieldDesc
	Returns []LuaReturnDesc
}

func (funcAnnotation *LuaFuncAnnotation) Type() string {
	return FUNC
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
				func(arg LuaFieldDesc) string {
					return arg.Name
				}), ", "),
		),
	)

	builder.WriteString("\n")

	return builder.String()
}

//
// Class annotations
//

type LuaClassAnnotation struct {
	Name   string
	Fields []LuaFieldDesc
}

func (classAnnotation *LuaClassAnnotation) Type() string {
	return CLASS
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

func (classAnnotation *LuaClassAnnotation) Field(name string, kind string, desc string) *LuaClassAnnotation {
	classAnnotation.Fields = append(classAnnotation.Fields, LuaFieldDesc{
		Name: name,
		Type: kind,
		Desc: desc,
	})
	return classAnnotation
}
