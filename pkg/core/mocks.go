package core

import (
	"fmt"
	"strings"

	"github.com/patrixr/glue/pkg/luatools"
	lua "github.com/yuin/gopher-lua"
)

func (glue *Glue) generateMockMethod(name string) lua.LGFunction {
	return func(st *lua.LState) int {
		inputs := luatools.GetAllArgsAsStrings(st)
		text := fmt.Sprintf("(mocked) %s(%s)", name, strings.Join(inputs, ", "))

		glue.Log.Info(text)

		return 0
	}
}
