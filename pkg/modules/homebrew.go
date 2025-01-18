package modules

import (
	"github.com/patrixr/glue/pkg/core"
	. "github.com/patrixr/glue/pkg/machine"
	. "github.com/patrixr/glue/pkg/runtime"
)

func init() {
	Registry.RegisterModule(HomebrewMod)
}

func HomebrewMod(glue *core.Glue) error {
	// @TODO: type the array's items (e.g string[])

	ensure := func(R Runtime, args *Arguments) (RTValue, error) {
		if !glue.Verbose {
			glue.Log.Quiet()
		}

		defer glue.Log.Loud()

		if err := InstallHomebrew(glue.Machine, glue.Log.Stdout, glue.Log.Stderr); err != nil {
			return nil, err
		}
		return nil, UpdateHomebrew(glue.Machine, glue.Log.Stdout, glue.Log.Stderr)
	}

	mainHomebrew := func(R Runtime, args *Arguments) (RTValue, error) {
		params, err := DecodeDict[HomebrewParams](args.EnsureDict(0))

		if err != nil {
			return nil, err
		}

		if !glue.Verbose {
			glue.Log.Quiet()
		}

		defer glue.Log.Loud()

		return nil, HomebrewBundle(glue.Machine, params, glue.Log.Stdout, glue.Log.Stderr)
	}

	upgrade := func(R Runtime, args *Arguments) (RTValue, error) {
		if !glue.Verbose {
			glue.Log.Quiet()
		}

		defer glue.Log.Loud()

		return nil, HomebrewUpgrade(glue.Machine, glue.Log.Stdout, glue.Log.Stderr)
	}

	glue.Plug("homebrew_install", core.MODULE).
		Brief("Installs Homebrew if not already installed").
		Do(ensure)

	StringArray := TypedArray(STRING)

	glue.Plug("homebrew", core.MODULE).
		Brief("Marks a homebrew package for installation").
		Arg("params", CustomStruct("HomebrewParams", []Field{
			NewField("packages?", StringArray, "the homebrew packages to install"),
			NewField("taps?", StringArray, "the homebrew taps to install"),
			NewField("mas?", StringArray, "the homebrew mac app stores to install"),
			NewField("whalebrews?", StringArray, "the whalebrews install"),
			NewField("casks?", StringArray, "the homebrew casks to install"),
		}), "the packages to install").
		Do(mainHomebrew)

	glue.Plug("homebrew_upgrade", core.MODULE).
		Brief("Upgrades all homebrew packages").
		Do(upgrade)

	return nil
}
