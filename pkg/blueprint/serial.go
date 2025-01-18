package blueprint

import "strings"

type SerialBlueprint struct {
	Name       string      `json:"name"`
	Details    string      `json:"details"`
	Annotation string      `json:"annotation"`
	Children   []Blueprint `json:"children"`
	Function   BlueprintFunc
}

func NewSerialBlueprint(name string) *SerialBlueprint {
	return &SerialBlueprint{
		Name:     name,
		Children: []Blueprint{},
	}
}

func (blueprint *SerialBlueprint) Execute() Results {
	results := Results{}

	if blueprint.Function != nil {
		trace := blueprint.Function()

		if trace.Error != nil {
			results.ErrorCount++
			results.Success = false
		}

		results.Traces = append(results.Traces, trace)
	}

	for _, child := range blueprint.Children {
		res := child.Execute()
		results.Traces = append(results.Traces, res.Traces...)
		results.ErrorCount += res.ErrorCount
		results.Success = results.Success && res.Success
	}

	return results
}

func (blueprint *SerialBlueprint) Action(name string, details string, usertext string, fn ActionFunc) {
	blueprint.Children = append(blueprint.Children, &SerialBlueprint{
		Function: func() Trace {
			err := fn()
			return Trace{
				Name:  name,
				Error: err,
			}
		},
		Name:       name,
		Details:    details,
		Annotation: usertext,
		Children:   []Blueprint{},
	})
}

func (blueprint *SerialBlueprint) Add(subBlueprint Blueprint) {
	blueprint.Children = append(blueprint.Children, subBlueprint)
}

func (blueprint *SerialBlueprint) PrettyPrint() string {
	builder := strings.Builder{}
	blueprint.prettyPrintRecursive(&builder, 0)
	return builder.String()
}

// (internal)
func (blueprint *SerialBlueprint) prettyPrintRecursive(builder *strings.Builder, depth int) {
	for i := 0; i < depth; i++ {
		builder.WriteString("  ")
	}
	builder.WriteString("+ ")
	builder.WriteString(blueprint.Name)
	builder.WriteString("\n")

	for _, child := range blueprint.Children {
		childBlueprint := child.(*SerialBlueprint)
		childBlueprint.prettyPrintRecursive(builder, depth+1)
	}
}
