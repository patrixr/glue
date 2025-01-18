package lua

import (
	"fmt"
	"strings"

	"github.com/patrixr/glue/pkg/core"
	"github.com/patrixr/glue/pkg/runtime"
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

func GenerateTypeDefinitions(glue *core.Glue) string {
	annotations := LuaAnnotations{}

	for _, mod := range glue.Modules {
		path := strings.Split(mod.Name, ".")
		annotations.tryAddNestedField(path[:len(path)-1])

		for _, arg := range mod.Args {
			customStructType, isCustomStruct := arg.Type.(runtime.CustomStructType)

			if !isCustomStruct {
				continue
			}

			annotations.AddClass(customStructType.Name(), customStructType.Fields)
		}

		annotations.AddFunction(mod)

	}

	return annotations.Render()
}

// func GenerateFunctionType(mod *core.GluePlugin) string {

// }

// func GenerateCustomStructDefinition(arg runtime.Type) string {

// }

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

func (annotations *LuaAnnotations) AddFunction(mod *core.GluePlugin) {
	annotations.items = append(annotations.items, &LuaFuncAnnotation{mod})
}

func (annotations *LuaAnnotations) AddClass(name string, fields []runtime.Field) {
	class := &LuaClassAnnotation{
		Name:   name,
		Fields: fields,
	}
	annotations.Add(class)
}

func (annotations *LuaAnnotations) Type() string {
	return META
}

func (annotations *LuaAnnotations) FindAllByType(kind string) []Annotable {
	return q.FindAll(annotations.items, func(a Annotable, _ int) bool {
		return a.Type() == kind
	})
}

func (annotations *LuaAnnotations) findTable(name string) *LuaTableAnnotation {
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

func (annotations *LuaAnnotations) tryAddNestedField(path []string) {
	if len(path) == 0 {
		return
	}

	head := path[0]
	tail := path[1:]
	base := annotations.findTable(head)

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

// ------------------------------------------
// Table annotations (for nested fields)
// ------------------------------------------

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

// ------------------------------------------
// Function annotations
// ------------------------------------------

type LuaFuncAnnotation struct {
	plug *core.GluePlugin
}

func (mod *LuaFuncAnnotation) Type() string {
	return FUNC
}

func (funcAnnotation *LuaFuncAnnotation) Render() string {
	var builder strings.Builder

	builder.WriteString("---\n")

	for _, line := range strings.Split(funcAnnotation.plug.Brief, "\n") {
		builder.WriteString(fmt.Sprintf("--- %s\n", line))
	}

	builder.WriteString("---\n")

	for _, arg := range funcAnnotation.plug.Args {
		builder.WriteString(fmt.Sprintf("---@param %s %s %s\n", arg.Name, runtime.TypeName(arg.Type), arg.Desc))
	}

	builder.WriteString("---\n")

	if ret := funcAnnotation.plug.ReturnType; ret != nil {
		name := ret.Name()
		desc := ""

		typeWithDesc, hasDesc := ret.(runtime.DescribableType)
		if hasDesc {
			desc = typeWithDesc.Describe()
		}

		builder.WriteString(fmt.Sprintf("---@return %s %s\n", name, desc))
	}

	builder.WriteString("---\n")

	builder.WriteString(
		fmt.Sprintf("function %s(%s) end",
			funcAnnotation.plug.Name,
			strings.Join(q.Map(funcAnnotation.plug.Args,
				func(arg runtime.ArgDef) string {
					return arg.Name
				}), ", "),
		),
	)

	builder.WriteString("\n")

	return builder.String()
}

// ------------------------------------------
// Class annotations
// ------------------------------------------

type LuaClassAnnotation struct {
	Name   string
	Fields []runtime.Field
}

func (classAnnotation *LuaClassAnnotation) Type() string {
	return CLASS
}

func (classAnnotation *LuaClassAnnotation) Render() string {
	var builder strings.Builder

	builder.WriteString(fmt.Sprintf("---@class %s\n", classAnnotation.Name))

	for _, field := range classAnnotation.Fields {
		optstr := ""
		if field.Optional {
			optstr = "?"
		}
		builder.WriteString(fmt.Sprintf("---@field %s%s %s %s\n", field.Name, optstr, runtime.TypeName(field.Type), field.Desc))
	}

	builder.WriteString("\n")

	return builder.String()
}
