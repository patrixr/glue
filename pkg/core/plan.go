package core

import "strings"

type StepFunc func() error

type Plan interface {
	Execute() []error
	Step(name string, fn StepFunc)
	Add(plan Plan)
	Pretty() string
}

type TreePlan struct {
	Name     string `json:"name"`
	Fn       StepFunc
	Children []Plan `json:"children"`
}

func (step *TreePlan) Execute() []error {
	var errors []error = []error{}

	if step.Fn != nil {
		if err := step.Fn(); err != nil {
			errors = append(errors, err)
		}
	}

	for _, step := range step.Children {
		errs := step.Execute()
		errors = append(errors, errs...)
	}

	return errors
}

func (step *TreePlan) Step(name string, fn StepFunc) {
	step.Children = append(step.Children, &TreePlan{
		Fn:       fn,
		Name:     name,
		Children: []Plan{},
	})
}

func (step *TreePlan) Add(plan Plan) {
	step.Children = append(step.Children, plan)
}

func NewPlan(name string) Plan {
	return &TreePlan{
		Fn:       nil,
		Name:     name,
		Children: []Plan{},
	}
}

func (step *TreePlan) Pretty() string {
	builder := strings.Builder{}
	step.prettyRecursive(&builder, 0)
	return builder.String()
}

// (internal)
func (step *TreePlan) prettyRecursive(builder *strings.Builder, depth int) {
	for i := 0; i < depth; i++ {
		builder.WriteString("  ")
	}
	builder.WriteString("+ ")
	builder.WriteString(step.Name)
	builder.WriteString("\n")

	for _, child := range step.Children {
		childPlan := child.(*TreePlan)
		childPlan.prettyRecursive(builder, depth+1)
	}
}
