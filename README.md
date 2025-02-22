# Glue

Glue is a powerful system setup tool combining using Lua as a configuration language .
It provides an intuitive, imperative approach to machine setup and configuration management.

> Warning: Glue is still is in prototype phase. Features and APIs may evolve

## Overview

Glue seamlessly combines Go + Lua (hence the name) to offer:

- Lua-based configuration language
- Simple, imperative API (inspired by Ansible)
- Extensible module system
- "Blueprint" generation and execution
- Filtering of configuration blocks

> **Note:** Glue is still in prototype phase. Features and APIs may evolve. Please report issues via GitHub.

## Quick Start

Glue typically works globally on your system by referencing the `glue.lua` file in your XDG_CONFIG_HOME folder.
Typically that would be `~/.config/glue/glue.lua`.

Here's an example of a configuration that sets up some configurations and installs Homebrew packages:

```lua
group("configs", function ()
    Copy({
        source = "./configs/alacritty" .. name,
        dest = "~/.config/alacritty",
        strategy = "merge"
    })

    Blockinfile({
        state = true,
        block = read("./configs/zshrc.sh"),
        path = "~/.zshrc"
    })
end)

group("homebrew", function ()
    HomebrewInstall()

    Homebrew({
        taps =  {
            "oven-sh/bun",
            "homebrew/cask-fonts",
        },
        casks = {
            "zen-browser",
            "steam",
            "emacs",
            "love",
        },
        packages = {
            "ffmpeg",
            "watch",
            "httpie",
            "ruby",
            "lua",
        }
    })
end)

```

## CLI Reference

```bash
glue [flags] [command]
```

### Commands

| Command      | Description                              |
| ------------ | ---------------------------------------- |
| `completion` | Generate shell autocompletion scripts    |
| `document`   | Generate internal function documentation |
| `help`       | Display help information                 |
| `init`       | Initialize Glue on your system           |
| `only`       | Execute specific groups using a selector |

### Flags

| Flag                | Description                                            |
| ------------------- | ------------------------------------------------------ |
| `--plan`            | See the execution blueprints without applying anything |
| `-h, --help`        | Show help information                                  |
| `-p, --path string` | Specify glue.lua location                              |
| `-v, --verbose`     | Enable verbose logging                                 |

## Extending Glue

Glue can be extended through its module system. Create new modules in the `modules` package using the registry system.

Here's an example of a simple module that prints a message:

```go
Registry.RegisterModule(func(glue *core.Glue) error {
	glue.Plug("print", core.FUNCTION).
		Brief("Print a string").
		Arg("obj", ANY, "the message or object to log").
		Do(func(R Runtime, args *Arguments) (RTValue, error) {
			glue.Log.Info(args.Get(0).String())
			return nil, nil
		})

	return nil
})
```

## Contributing

Contributions are welcome! Please:

1. Fork the repository
2. Create a feature branch
3. Submit a pull request

## License

GNU General Public License v3.0
