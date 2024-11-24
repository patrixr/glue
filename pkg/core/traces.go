package core

import (
	"github.com/patrixr/glue/pkg/luatools"
	lua "github.com/yuin/gopher-lua"
)

type Trace struct {
	Name  string
	Args  []string
	Error error
	About string
}

func (glue *Glue) SaveTrace(name string, L *lua.LState, err error) {
	glue.ExecutionTrace = append(glue.ExecutionTrace, Trace{
		Name:  name,
		Args:  luatools.GetAllArgsAsStrings(L),
		Error: err,
	})

	glue.Fire(EV_NEW_TRACE, &glue.ExecutionTrace[len(glue.ExecutionTrace)-1])
}
