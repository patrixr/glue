package modules

import (
	"github.com/patrixr/glue/pkg/core"
	"github.com/patrixr/glue/pkg/homebrew"
	"github.com/patrixr/glue/pkg/luatools"
)

func init() {
	Registry.RegisterModule(HomebrewMod)
}

func HomebrewMod(glue *core.Glue) error {
	brew := homebrew.NewHomebrew()

	pkg := luatools.StrFunc(func(name string) error {
		brew.Brew(name)
		return nil
	})

	cask := luatools.StrFunc(func(name string) error {
		brew.Cask(name)
		return nil
	})

	tap := luatools.StrFunc(func(name string) error {
		brew.Tap(name)
		return nil
	})

	mas := luatools.StrFunc(func(name string) error {
		brew.Mas(name)
		return nil
	})

	whalebrew := luatools.StrFunc(func(name string) error {
		brew.Whalebrew(name)
		return nil
	})

	sync := luatools.Func(func() error {
		return brew.Install()
	})

	glue.AddFunction("brew.package", pkg)
	glue.AddFunction("brew.cask", cask)
	glue.AddFunction("brew.tap", tap)
	glue.AddFunction("brew.mas", mas)
	glue.AddFunction("brew.whalebrew", whalebrew)
	glue.AddFunction("brew.sync", sync)

	return nil
}
