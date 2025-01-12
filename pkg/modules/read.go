package modules

import (
	"os"

	"github.com/patrixr/glue/pkg/core"
	"github.com/patrixr/glue/pkg/luatools"
)

func init() {
	Registry.RegisterModule(
		func(glue *core.Glue) error {
			glue.Plug().
				Name("read").
				Short("Reads a file as a string").
				Long("Reads a file as a string").
				Arg("path", "string", "the path of the file to read").
				Return("string", "the file content").
				Example("read(\"./some/file\")").
				Mode(core.READ).
				MockReturn(luatools.EmptyLuaString).
				Do(luatools.StrInStrOutFunc(func(path string) (string, error) {
					resolvedPath, err := glue.SmartPath(path)

					if err != nil {
						return "", err
					}

					data, err := os.ReadFile(resolvedPath)

					if err != nil {
						return "", err
					}

					return string(data), err
				}))

			return nil
		})
}
