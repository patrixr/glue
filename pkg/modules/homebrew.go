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
	brew := homebrew.NewHomebrew(glue.Log.Stdout, glue.Log.Stderr)

	ensure := luatools.Func(func() error {
		if !glue.Verbose {
			glue.Log.Quiet()
		}

		defer glue.Log.Loud()

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
		if !glue.Verbose {
			glue.Log.Quiet()
		}

		defer glue.Log.Loud()

		return brew.Install()
	})

	upgrade := luatools.Func(func() error {
		if !glue.Verbose {
			glue.Log.Quiet()
		}

		defer glue.Log.Loud()

		return brew.Upgrade()
	})

	glue.Plug().
		Name("homebrew_install").
		Short("Installs Homebrew if not already installed").
		Example("homebrew_install()").
		Do(ensure)

	glue.Plug().
		Name("homebrew").
		Short("Marks a homebrew package for installation").
		Arg("pkg", "string", "the name of the package to install").
		Example("homebrew('git')").
		Example("homebrew('zsh)").
		Example("homebrew(\"package\")").
		Example("homebrew_sync()").
		Do(pkg)

	glue.Plug().
		Name("homebrew_cask").
		Short("Marks a cask for installation").
		Arg("pkg", "string", "the name of the cask to install").
		Example("homebrew_cask('firefox')").
		Example("homebrew_cask('spotify')").
		Example("homebrew_sync()").
		Do(cask)

	glue.Plug().
		Name("homebrew_tap").
		Short("Marks a homebrew tap for installation").
		Arg("tap", "string", "the name of the tap to install").
		Example("homebrew_tap('homebrew/cask')").
		Do(tap)

	glue.Plug().
		Name("homebrew_mas").
		Short("Marks a Mac App Store package for installation").
		Arg("name", "string", "the name of the mas to install").
		Example("homebrew_mas('1Password')").
		Example("homebrew_mas('Slack')").
		Example("homebrew_sync()").
		Do(mas)

	glue.Plug().
		Name("homebrew_whalebrew").
		Short("Marks a whalebrew package for installation").
		Arg("name", "string", "the name of the whalebrew to install").
		Example("homebrew_whalebrew('whalebrew/awscli'").
		Example("homebrew_whalebrew('whalebrew/ffmpeg')").
		Example("homebrew_sync()").
		Do(whalebrew)

	glue.Plug().
		Name("homebrew_sync").
		Short("Installs all marked packages").
		Example("homebrew_sync()").
		Do(sync)

	glue.Plug().
		Name("homebrew_upgrade").
		Short("Upgrades all homebrew packages").
		Example("homebrew_upgrade()").
		Do(upgrade)

	return nil
}
