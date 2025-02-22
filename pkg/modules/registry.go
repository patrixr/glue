package modules

import (
	"github.com/patrixr/glue/pkg/core"
)

// @auteur("Modules")
//
// # Modules
//
// Modules in Glue are self-contained units of functionality that can be independently developed and integrated into the Glue framework.
//
// In practice they look like making a Lua function call, such as a `Copy(...)` call to copy files over, but under the hood they actually
// get composed into a **Blueprint** which can be viewed, serialized and shared with or without being executed.
//
// Modules are **always capitalized** (e.g. `Blockinfile`)

// @auteur("Helpers")
//
// # Helper functions
//
// Helper functions look similar to modules in the sense that they are also Lua functions.
// The key difference is that helper functions get **executed immediatly** during the Blueprint creation, and are *not* part of the generated blueprints.
//
// Some common use cases involve reading from files, string manipulation or http requests.
//
// Helpers are **always lowercased** (e.g. `read`)
const AUTEUR_MODUlES = "Modules"
const AUTEUR_UTILS = "Utils"

type ModuleInstaller func(glue *core.Glue) error

type ModuleRegistry struct {
	installers []ModuleInstaller
}

func NewModuleRegistry() *ModuleRegistry {
	return &ModuleRegistry{}
}

var Registry *ModuleRegistry = NewModuleRegistry()

// RegisterModule adds a module installer to the registry
// The installer is a function that installs the modules onto the Glue instance
func (r *ModuleRegistry) RegisterModule(mod ModuleInstaller) {
	r.installers = append(r.installers, mod)
}

// InstallModules installs all the modules in the registry onto the Glue instance
func (r *ModuleRegistry) InstallModules(glue *core.Glue) error {
	for _, install := range r.installers {
		if err := install(glue); err != nil {
			return err
		}
	}
	return nil
}
