package modules

import (
	"os"
	"strings"

	"github.com/patrixr/glue/pkg/core"
	"github.com/patrixr/glue/pkg/luatools"
	"github.com/patrixr/glue/pkg/shell"
	"github.com/patrixr/q"
	lua "github.com/yuin/gopher-lua"
)

func init() {
	Registry.RegisterModule(UtilitiesMod)
}

func UtilitiesMod(glue *core.Glue) error {
	print := func(L *lua.LState) (int, error) {
		input := luatools.GetArgAsString(L, 1)
		glue.Log.Info(input)
		return 0, nil
	}

	sh := luatools.StrFunc(func(cmd string) error {
		return shell.Run(cmd, os.Stdout, os.Stderr)
	})

	trim := luatools.StrInStrOutFunc(func(s string) (string, error) {
		return strings.TrimSpace(q.TrimIndent(s)), nil
	})

	read := luatools.StrInStrOutFunc(func(path string) (string, error) {
		resolvedPath, err := glue.SmartPath(path)

		if err != nil {
			return "", err
		}

		data, err := os.ReadFile(resolvedPath)

		if err != nil {
			return "", err
		}

		return string(data), err
	})

	glue.Plug().
		Name("sh").
		Short("Run a shell command").
		Long("Run a shell command").
		Arg("cmd", "string", "the shell command to run").
		Example("sh('ls')").
		Do(sh)

	glue.Plug().
		Name("print").
		Short("Print a string").
		Long("Print a string").
		Arg("obj", "any", "the message or object to log").
		Example("print('Hello, world!')").
		Do(print)

	glue.Plug().
		Name("trim").
		Short("Trims the extra indentation of a multi-line string").
		Long("Trims the extra indentation of a multi-line string").
		Arg("txt", "string", "the text to trim").
		Return("string", "the trimmed text").
		Example("trim(text)").
		Mode(core.NONE).
		Bypass().
		Do(trim)

	glue.Plug().
		Name("read").
		Short("Reads a file as a string").
		Long("Reads a file as a string").
		Arg("path", "string", "the path of the file to read").
		Return("string", "the file content").
		Example("read(\"./some/file\")").
		Mode(core.READ).
		MockReturn(luatools.EmptyLuaString).
		Do(read)

	return nil
}
