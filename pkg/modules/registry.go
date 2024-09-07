package modules

import (
	"github.com/patrixr/glue/pkg/core"
)

type ModuleInstaller func(glue *core.Glue) error

type ModuleRegistry struct {
	modules []ModuleInstaller
}

func NewModuleRegistry() *ModuleRegistry {
	return &ModuleRegistry{}
}

var Registry *ModuleRegistry = NewModuleRegistry()

// Functions

func (r *ModuleRegistry) RegisterModule(mod ModuleInstaller) {
	r.modules = append(r.modules, mod)
}

func (r *ModuleRegistry) InstallModules(glue *core.Glue) error {
	for _, mod := range r.modules {
		if err := mod(glue); err != nil {
			return err
		}
	}
	return nil
}
