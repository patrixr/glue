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

	ensure := luatools.Func(func() error {
		if err := homebrew.InstallHomebrew(glue.Log.Stdout, glue.Log.Stderr); err != nil {
			return err
		}
		return homebrew.UpdateHomebrew(glue.Log.Stdout, glue.Log.Stderr)
	})

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

	glue.Plug().
		Name("brew.ensure").
		Short("Installs Homebrew if not already installed").
		Example("brew.ensure()").
		Do(ensure)

	glue.Plug().
		Name("brew.package").
		Short("Marks a homebrew package for installation").
		Example("brew.package('git')").
		Example("brew.package('zsh)").
		Example("brew.sync()").
		Do(pkg)

	glue.Plug().
		Name("brew.cask").
		Short("Marks a cask for installation").
		Example("brew.cask('firefox')").
		Example("brew.cask('spotify')").
		Example("brew.sync()").
		Do(cask)

	glue.Plug().
		Name("brew.tap").
		Short("Marks a homebrew tap for installation").
		Example("brew.tap('homebrew/cask')").
		Do(tap)

	glue.Plug().
		Name("brew.mas").
		Short("Marks a Mac App Store package for installation").
		Example("brew.mas('1Password')").
		Example("brew.mas('Slack')").
		Example("brew.sync()").
		Do(mas)

	glue.Plug().
		Name("brew.whalebrew").
		Short("Marks a whalebrew package for installation").
		Example("brew.whalebrew('whalebrew/awscli'").
		Example("brew.whalebrew('whalebrew/ffmpeg')").
		Do(whalebrew)

	glue.Plug().
		Name("brew.sync").
		Short("Installs all marked packages").
		Example("brew.sync()").
		Do(sync)

	return nil
}
