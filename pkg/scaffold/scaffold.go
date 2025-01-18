package scaffold

import (
	"github.com/patrixr/auteur"
)

func init() {
	auteur.Write("Concepts", `
		Scaffolding
		============

		Glue provides way of scaffolding a folder on your behalf.
		This is key to getting the LSP to work properly as it will generate the necessary files for it to work.
	`)
}

type Scaffold interface {
	Setup(folder string) error
	Typegen() string
}
