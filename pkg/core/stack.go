package core

import (
	"github.com/patrixr/q"
)

type ScriptType int

const (
	FILE ScriptType = 1 << iota
	STRING
	REMOTE
)

type GlueCodeGroup struct {
	Name        string
	Annotations map[string]string
}

type GlueScript struct {
	Uri        string
	Type       ScriptType
	GroupStack []*GlueCodeGroup
	GroupNames []string
}

type GlueStack struct {
	ExecutionStack []*GlueScript
}

func NewGlueScope(fileName string) *GlueStack {
	return &GlueStack{
		ExecutionStack: []*GlueScript{
			{
				Uri: fileName,
				GroupStack: []*GlueCodeGroup{
					{Name: RootLevel},
				},
			},
		},
	}
}

func (scope *GlueStack) PushScript(file string, kind ScriptType) {
	scope.ExecutionStack = append(scope.ExecutionStack, &GlueScript{
		Uri:  file,
		Type: kind,
		GroupStack: []*GlueCodeGroup{
			{
				Name:        RootLevel,
				Annotations: map[string]string{},
			},
		},
	})
}

func (scope *GlueStack) PopScript() {
	scope.ExecutionStack = scope.ExecutionStack[:len(scope.ExecutionStack)-1]
}

func (scope *GlueStack) HasActiveScript() bool {
	return len(scope.ExecutionStack) > 0
}

func (scope *GlueStack) ActiveScript() *GlueScript {
	q.Assert(len(scope.ExecutionStack) > 0, "Accessing file on an empty file stack")
	return scope.ExecutionStack[len(scope.ExecutionStack)-1]
}

func (scope *GlueStack) PushGroup(name string) {
	executable := scope.ActiveScript()
	executable.GroupStack = append(executable.GroupStack, &GlueCodeGroup{
		Name:        name,
		Annotations: map[string]string{},
	})
}

func (scope *GlueStack) PopGroup() {
	executable := scope.ActiveScript()
	q.Assert(len(executable.GroupStack) > 0, "Trying to pop a group on an empty stack")
	executable.GroupStack = executable.GroupStack[:len(executable.GroupStack)-1]
}

func (scope *GlueStack) CurrentGroup() *GlueCodeGroup {
	script := scope.ActiveScript()
	q.Assert(len(script.GroupStack) > 0, "Trying to pop a group on an empty stack")
	return script.GroupStack[len(script.GroupStack)-1]
}

func (scope *GlueStack) AnnotateCurrentGroup(key string, value string) {
	group := scope.CurrentGroup()
	group.Annotations[key] = value
}

func (scope *GlueStack) GetActiveGroupAnnotation(key string) string {
	group := scope.CurrentGroup()
	val, ok := group.Annotations[key]

	if !ok {
		return ""
	}
	return val
}

func (grp *GlueCodeGroup) Set(key string, value string) {
	grp.Annotations[key] = value
}

func (grp *GlueCodeGroup) Get(key string) (string, bool) {
	val, ok := grp.Annotations[key]
	return val, ok
}
