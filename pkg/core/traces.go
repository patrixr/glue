package core

type Trace struct {
	Name    string
	Details string
	Error   error
	About   string
}

func (glue *Glue) SaveTrace(name string, details string, err error) {
	glue.ExecutionTrace = append(glue.ExecutionTrace, Trace{
		Name:    name,
		Details: details,
		Error:   err,
	})

	glue.Fire(EV_NEW_TRACE, &glue.ExecutionTrace[len(glue.ExecutionTrace)-1])
}
