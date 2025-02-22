---
path: /getting-started
priority: 10
---

# Getting Started

## Installation

To install Auteur, you can use one of the following methods:

### Homebrew

Glue is available as a Homebrew formula. To install Auteur using Homebrew, run the following commands:

First to set up the tap

```bash
brew tap tronica/tap
```

And then to install Glue

```bash
brew install glue
```

Glue should then be available as a command line tool.

### Using Go

To install Glue using Go, run the following command:

```bash
go install github.com/patrixr/glue
```

## Initialization

To initialize a new Glue project, run the following command:

```bash
glue init
```

This should create a new directory according to your `XDG_CONFIG_HOME` or `HOME` environment variable.

Typically, this will be `~/.config/glue`

## Editing the script

To edit the script, navigate to the newly created directory and open the `glue.lua` file.

Running glue will always look for a `glue.lua` file in the current directory first, and then default to the `~/.config/glue` directory.
