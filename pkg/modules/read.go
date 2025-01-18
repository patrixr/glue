package modules

import (
	"os"

	"github.com/patrixr/glue/pkg/core"
	. "github.com/patrixr/glue/pkg/runtime"
)

func init() {
	Registry.RegisterModule(
		func(glue *core.Glue) error {
			glue.Plug("read", core.FUNCTION).
				Brief("Reads a file as a string").
				Arg("path", STRING, "the path of the file to read").
				Return(STRING, "the file content").
				Do(func(R Runtime, args *Arguments) (RTValue, error) {
					path := args.EnsureString(0).String()
					resolvedPath, err := glue.SmartPath(path)

					if err != nil {
						return nil, err
					}

					data, err := os.ReadFile(resolvedPath)

					if err != nil {
						return nil, err
					}

					return R.String(string(data)), nil
				})

			return nil
		})
}
