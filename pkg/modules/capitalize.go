package modules

import (
	"strings"

	"github.com/patrixr/glue/pkg/core"
	. "github.com/patrixr/glue/pkg/runtime"
)

func init() {
	// @auteur("Helpers/Capitalize")
	//
	// # Capitalize
	//
	// This module provides a function to capitalize the first letter of a given string.
	// It takes a single string argument and returns the string with its first letter converted to uppercase.
	// This can be useful for formatting text inputs where proper capitalization is required
	//
	// Example:
	//
	// ```lua
	// capitalize("hello world") // returns "Hello world"
	// ```
	Registry.RegisterModule(func(glue *core.Glue) error {
		glue.Plug("capitalize", core.FUNCTION).
			Brief("Uppercase the first letter of a string").
			Arg("txt", STRING, "the text to capitalize").
			Return(STRING, "the text with capitalized first letter").
			Do(func(R Runtime, args *Arguments) (RTValue, error) {
				s := args.EnsureString(0).String()
				return R.String(strings.ToUpper(s[:1]) + s[1:]), nil
			})

		return nil
	})
}
