// @auteur("Concepts")
//
// # Scaffolding
//
// Glue provides way of scaffolding a folder on your behalf.
// This is key to getting the LSP to work properly as it will generate the necessary files for it to work.
// In the case of Lua, this means creating (or updating) the project's .luarc file to include Glue's definitions
package scaffold

type Scaffold interface {
	Setup(folder string) error
	Typegen() string
}
