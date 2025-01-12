package modules

import (
	"github.com/patrixr/glue/pkg/core"
	"github.com/patrixr/glue/pkg/homebrew"
	"github.com/patrixr/glue/pkg/luatools"
)

func init() {
	Registry.RegisterModule(HomebrewMod)
}

type HomebrewParams struct {
	Packages   []string `json:"packages"`
	Casks      []string `json:"casks"`
	Taps       []string `json:"taps"`
	Mas        []string `json:"mas"`
	Whalebrews []string `json:"whalebrews"`
}

func HomebrewMod(glue *core.Glue) error {
	glue.Annotations.AddClass("HomebrewParams").
		Field("packages?", "string[]", "the homebrew packages to install").
		Field("taps?", "string[]", "the homebrew taps to install").
		Field("mas?", "string[]", "the homebrew mac app stores to install").
		Field("whalebrews?", "string[]", "the whalebrews install").
		Field("casks?", "string[]", "the homebrew casks to install")

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

	mainHomebrew := luatools.TableFunc(func(params HomebrewParams) error {
		brew := homebrew.NewHomebrew(glue.Log.Stdout, glue.Log.Stderr)

		if !glue.Verbose {
			glue.Log.Quiet()
		}

		defer glue.Log.Loud()

		for _, it := range params.Packages {
			brew.Package(it)
		}
		for _, it := range params.Casks {
			brew.Cask(it)
		}
		for _, it := range params.Taps {
			brew.Tap(it)
		}
		for _, it := range params.Mas {
			brew.Mas(it)
		}
		for _, it := range params.Whalebrews {
			brew.Whalebrew(it)
		}

		return brew.Install()
	})

	upgrade := luatools.Func(func() error {
		if !glue.Verbose {
			glue.Log.Quiet()
		}

		defer glue.Log.Loud()

		brew := homebrew.NewHomebrew(glue.Log.Stdout, glue.Log.Stderr)

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
		Arg("params", "HomebrewParams", "the packages to install").
		Do(mainHomebrew)

	glue.Plug().
		Name("homebrew_upgrade").
		Short("Upgrades all homebrew packages").
		Example("homebrew_upgrade()").
		Do(upgrade)

	return nil
}
