package core

func WithDryRun(dryRun bool) GlueOption {
	return func(glue *Glue) {
		glue.DryRun = dryRun
	}
}
