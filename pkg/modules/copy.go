package modules

import (
	"errors"
	"fmt"
	"os"

	cp "github.com/otiai10/copy"
	"github.com/patrixr/glue/pkg/core"
	"github.com/patrixr/glue/pkg/luatools"
)

func init() {
	Registry.RegisterModule(func(glue *core.Glue) error {
		glue.Annotations.AddClass("CopyOpts").
			Field("source", "string", "the file or folder to copy").
			Field("dest", "string", "the destination to copy to").
			Field("strategy?", `"replace"|"merge"`, "a strategy for how to manage conflicts (defaults to merge)").
			Field("symlink?", `"deep"|"shallow"|"skip"`, "how to handle symlinks (copy the content, copy the link, or the default skip)")

		glue.Plug().
			Name("copy").
			Short("Copies folder").
			Long("Copies").
			Arg("opts", "CopyOpts", "the copy options").
			Do(luatools.TableFunc[CopyOpts](func(opts CopyOpts) error {
				dest, err := glue.SmartPath(opts.Dest)

				if err != nil {
					return err
				}

				src, err := glue.SmartPath(opts.Source)

				if err != nil {
					return err
				}

				glue.Log.Info("[Copy]", "src", opts.Source, "dst", opts.Dest)

				opts.Dest = dest
				opts.Source = src

				return Copy(opts)
			}))

		return nil
	})
}

const (
	StrategyMerge   = "merge"
	StrategyReplace = "replace"
	SymlinkDeep     = "deep"
	SymlinkShallow  = "shallow"
	SymlinkSkip     = "skip"
)

type CopyOpts struct {
	Strategy      string `json:"strategy"`
	Symlink       string `json:"symlink"`
	PreserveOwner bool   `json:"preserve_owner"`
	Source        string `json:"source"`
	Dest          string `json:"dest"`
}

func Copy(opts CopyOpts) error {
	src := opts.Source
	dst := opts.Dest

	if _, err := os.Stat(src); err != nil {
		return err
	}

	if len(dst) == 0 {
		return errors.New(fmt.Sprintf("Invalid copy destination %s", dst))
	}

	return cp.Copy(src, dst, cp.Options{
		PreserveOwner: true,
		OnDirExists: func(_ string, _ string) cp.DirExistsAction {
			if opts.Strategy == StrategyReplace {
				return cp.Replace
			}
			return cp.Merge
		},
		OnSymlink: func(_ string) cp.SymlinkAction {
			if opts.Symlink == SymlinkDeep {
				return cp.Deep
			}
			if opts.Symlink == SymlinkShallow {
				return cp.Shallow
			}
			return cp.Skip
		},
	})
}
