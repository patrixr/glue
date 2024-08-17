package modules

import (
	lua "github.com/yuin/gopher-lua"
)

type LuaLifecycle interface {
	AfterScript(func() error)
}

type LuaModuleLoader func(*lua.LState, LuaLifecycle) error

type ModuleRegistry struct {
	loaders []LuaModuleLoader
}

func NewModuleRegistry() *ModuleRegistry {
	return &ModuleRegistry{}
}

var Registry *ModuleRegistry = NewModuleRegistry()

func (r *ModuleRegistry) AddLoader(module LuaModuleLoader) {
	r.loaders = append(r.loaders, module)
}

func (r *ModuleRegistry) LoadModules(L *lua.LState, lifecycle LuaLifecycle) error {
	for _, module := range r.loaders {
		if err := module(L, lifecycle); err != nil {
			return err
		}
	}
	return nil
}
