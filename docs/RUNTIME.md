---
path: modules-vs-helpers
---

# Modules vs Helpers

The difference between Glue modules and helpers appears subtle at first, but it is important to understand the distinction.

## Modules

Modules appear as normal Lua functions, but they do not execute logic directly. Instead what they do is _register_ functions on the blueprint for it to be executed later.

Example:

```lua
Blockinfile({
  state = true,
  block = "some file content"
  path = "~/.zshrc"
})
```

This example script will populate a blueprint with a _Blockinfile_ step.
Notice that the module starts with a capital letter, which is a convention to indicate that it is a module.
Now let's expand this example by using a helper function.

## Helpers

Helpers are functions that execute logic directly. They are simple Lua methods.
For example, we'll use the `read` helper to dynamically load the content of the the _Blockinfile_ module into the blueprint.

```lua
Blockinfile({
  state = true,
  block = read("./local_config.zshrc"),
  path = "~/.zshrc"
})
```

This example script will read the content of the file `local_config.zshrc` and use it as the block content for the _Blockinfile_ module.
But that is executed during the blueprint generation, not during the execution of the blueprint.

Meaning that the content of the file will be serialized into the blueprint itself.
This is a critical concept to understand, as it would theoretically allow us to run the blueprint on a different machine without the need for the file to be present.

By convention, helper functions start with a lowercase letter.
