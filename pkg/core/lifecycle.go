package core

type GlueLifecycle struct {
	afterScriptFuncs  []func() error
	beforeScriptFuncs []func() error
}

func (l *GlueLifecycle) AfterScript(f func() error) {
	l.afterScriptFuncs = append(l.afterScriptFuncs, f)
}

func (l *GlueLifecycle) BeforeScript(f func() error) {
	l.beforeScriptFuncs = append(l.beforeScriptFuncs, f)
}

func (l *GlueLifecycle) RunBeforeScript() error {
	return runAll(l.beforeScriptFuncs)
}

func (l *GlueLifecycle) RunAfterScript() error {
	return runAll(l.afterScriptFuncs)
}

func NewGlueLifecycle() *GlueLifecycle {
	return &GlueLifecycle{}
}

func runAll(funcs []func() error) error {
	for _, f := range funcs {
		if err := f(); err != nil {
			return err
		}
	}
	return nil
}
