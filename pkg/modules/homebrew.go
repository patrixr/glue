package modules

import (
	"github.com/patrixr/glue/pkg/core"
	"github.com/patrixr/glue/pkg/homebrew"
	. "github.com/patrixr/glue/pkg/runtime"
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
	// @TODO: type the array's items (e.g string[])
	glue.Annotations.AddClass("HomebrewParams").
		Field("packages?", ARRAY, "the homebrew packages to install").
		Field("taps?", ARRAY, "the homebrew taps to install").
		Field("mas?", ARRAY, "the homebrew mac app stores to install").
		Field("whalebrews?", ARRAY, "the whalebrews install").
		Field("casks?", ARRAY, "the homebrew casks to install")

	ensure := func(R Runtime, args *Arguments) (RTValue, error) {
		if !glue.Verbose {
			glue.Log.Quiet()
		}

		defer glue.Log.Loud()

		if err := homebrew.InstallHomebrew(glue.Log.Stdout, glue.Log.Stderr); err != nil {
			return nil, err
		}
		return nil, homebrew.UpdateHomebrew(glue.Log.Stdout, glue.Log.Stderr)
	}

	mainHomebrew := func(R Runtime, args *Arguments) (RTValue, error) {
		params, err := DecodeDict[HomebrewParams](args.EnsureDict(0))

		if err != nil {
			return nil, err
		}

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

		return nil, brew.Install()
	}

	upgrade := func(R Runtime, args *Arguments) (RTValue, error) {
		if !glue.Verbose {
			glue.Log.Quiet()
		}

		defer glue.Log.Loud()

		brew := homebrew.NewHomebrew(glue.Log.Stdout, glue.Log.Stderr)

		return nil, brew.Upgrade()
	}

	glue.Plug().
		Name("homebrew_install").
		Short("Installs Homebrew if not already installed").
		Example("homebrew_install()").
		Do(ensure)

	// @TODO: type the class HomebrewParams
	glue.Plug().
		Name("homebrew").
		Short("Marks a homebrew package for installation").
		Arg("params", DICT, "the packages to install").
		Do(mainHomebrew)

	glue.Plug().
		Name("homebrew_upgrade").
		Short("Upgrades all homebrew packages").
		Example("homebrew_upgrade()").
		Do(upgrade)

	return nil
}
