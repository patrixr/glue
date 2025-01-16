package core

type StepFunc func() error

type Plan interface {
	Execute() []error
	Step(name string, fn StepFunc)
	Add(plan Plan)
}

type TreePlan struct {
	Name     string
	Fn       StepFunc
	Children []Plan
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
