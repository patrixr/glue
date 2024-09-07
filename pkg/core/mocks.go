package core

import (
	"fmt"
	"strings"

	lua "github.com/yuin/gopher-lua"
)

func (glue *Glue) generateMockMethod(name string) *lua.LFunction {
	L := glue.lstate

	return L.NewFunction(func(st *lua.LState) int {
		inputs := []string{}
		idx := 1

		for {
			arg := L.Get(idx)

			if arg == lua.LNil {
				break
			}

			inputs = append(inputs, L.ToStringMeta(arg).String())
			idx++
		}

		text := fmt.Sprintf("(mocked) %s(%s)", name, strings.Join(inputs, ", "))

		glue.Log.Info(text)

		return 0
	})
}
