---
path: runtime
priority: 9
---

# Lua Runtime

Glue is built on top of the Lua runtime, which provides a powerful and flexible scripting environment.

We've configured the runtime to act as a bit of a sandbox, allowing you to run scripts in a controlled environment.
Many of the Lua standard libraries were disabled for this purpose, but we left those that we believed would be harmless in the context of Glue.

The available libraries in this sandboxed environment include:

- **base**: This library provides the core functions of Lua, such as basic input and output, type checking, and other fundamental operations.
- **math**: This library includes mathematical functions, such as trigonometric functions, logarithms, and other common mathematical operations.
- **string**: This library offers functions for string manipulation, including pattern matching, finding substrings, and formatting.
- **table**: This library provides functions to manipulate tables, which are the primary data structure in Lua, allowing for operations like sorting, inserting, and removing elements.
