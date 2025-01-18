package modules

import (
	"github.com/patrixr/glue/pkg/core"
)

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
