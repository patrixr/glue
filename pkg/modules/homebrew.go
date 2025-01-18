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
	ensure := func(R Runtime, args *Arguments) (RTValue, error) {
		if !glue.Verbose {
			glue.Log.Quiet()
		}

		defer glue.Log.Loud()

		glue.Log.Info("Ensuring Homebrew is installed")

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

		glue.Log.Info("Running Homebrew")

		return nil, HomebrewBundle(glue.Machine, params, glue.Log.Stdout, glue.Log.Stderr)
	}

	upgrade := func(R Runtime, args *Arguments) (RTValue, error) {
		if !glue.Verbose {
			glue.Log.Quiet()
		}

		defer glue.Log.Loud()

		glue.Log.Info("Upgrading Homebrew packages")

		return nil, HomebrewUpgrade(glue.Machine, glue.Log.Stdout, glue.Log.Stderr)
	}

	glue.Plug("HomebrewInstall", core.MODULE).
		Brief("Installs Homebrew if not already installed").
		Do(ensure)

	StringArray := TypedArray(STRING)

	glue.Plug("Homebrew", core.MODULE).
		Brief("Marks a homebrew package for installation").
		Arg("params", CustomStruct("HomebrewParams", []Field{
			NewField("packages?", StringArray, "the homebrew packages to install"),
			NewField("taps?", StringArray, "the homebrew taps to install"),
			NewField("mas?", StringArray, "the homebrew mac app stores to install"),
			NewField("whalebrews?", StringArray, "the whalebrews install"),
			NewField("casks?", StringArray, "the homebrew casks to install"),
		}), "the packages to install").
		Do(mainHomebrew)

	glue.Plug("HomebrewUpgrade", core.MODULE).
		Brief("Upgrades all homebrew packages").
		Do(upgrade)

	return nil
}
