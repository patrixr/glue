# Glue

Glue allows machine setups using Lua as a configuration language. It is designed to be simple and easy to use, while still being powerful enough to handle complex configurations.

It is built in Golang (Go + Lua = Glue).

> **Warning:** This is an early release version of Glue. Features and functionality may change, and there may be bugs or incomplete features.
> Use with caution and report any issues you encounter.

It is inspired by Ansible's module system, but an imperative API is used instead of a declarative one.

Example:

```lua
note("Inject zsh configuration")

blockinfile({
  state = true,
  block = read("./my-zsh-config"),
  path = "~/.zshrc"
})

note("Import emacs configuration")

copy({
  source = "./my-emacs",
  dest = "~/.config/emacs",
  strategy = "merge"
})

group("packages", function ()
  homebrew_install()

  homebrew({
    casks = {
      "emacs",
      "love",
      "redisinsight",
      "docker",
      "obsidian",
      "ghostty",
      "firefox",
    },
    packages = {
      "typst",
      "gleam",
      "ruby",
      "go",
      "bun",
      "lua",
    }
  })
end)
```

## CLI Usage

```bash
Glue is a machine configuration tool that allows you to use Lua to easily streamline your system setup

Usage:
  glue [flags]
  glue [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  document    Generates a documentation of Glue's internal functions
  help        Help about any command
  init        Initializes glue on your system
  only        Execute Glue only on a subset of groups using a selector

Flags:
  -d, --dry-run       See the execution flow without running anything
  -h, --help          help for glue
  -p, --path string   Directory or file to look for glue.lua
  -v, --verbose       Enable verbose mode

Use "glue [command] --help" for more information about a command.
```

## Adding modules

Glue supports adding modules to extend its functionality. To add a module, create a new file in the `modules` package.
The file should register a module using the existing registry system.

```go
// Example of registering a print method
func init() {
	Registry.RegisterModule(func(glue *core.Glue) error {
		glue.Plug().
			Name("print").
			Short("Print a string").
			Long("Print a string").
			Arg("obj", "any", "the message or object to log").
			Example("print('Hello, world!')").
			Do(func(L *lua.LState) (int, error) {
				input := luatools.GetArgAsString(L, 1)
				glue.Log.Info(input)
				return 0, nil
			})

		return nil
	})
}
```
