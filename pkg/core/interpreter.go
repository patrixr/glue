package core

import (
	"github.com/patrixr/glue/pkg/modules"
	lua "github.com/yuin/gopher-lua"
)

func RunScriptFile(file string, modulesRegistry *modules.ModuleRegistry) error {
	lifecycle := NewGlueLifecycle()
	L := lua.NewState()
	defer L.Close()
	modulesRegistry.LoadModules(L, lifecycle)

	if err := L.DoFile(file); err != nil {
		return err
	}

	return lifecycle.RunAfterScript()
}

func AutoRun() error {
	file, err := AutoDetectScriptFile()

	if err != nil {
		return err
	}

	return RunScriptFile(file, modules.Registry)
}
