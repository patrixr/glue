package modules

import (
	"errors"
	"fmt"
	"os"

	cp "github.com/otiai10/copy"
	"github.com/patrixr/glue/pkg/core"
	. "github.com/patrixr/glue/pkg/runtime"
)

func init() {
	// @auteur("Modules/Copy")
	//
	// # Backup
	//
	// This module provides functionality to copy files or directories from a source to a destination.
	// It supports various options for handling conflicts and symlinks.
	//
	// ## Options
	// - `source`: The file or folder to copy.
	// - `dest`: The destination to copy to.
	// - `strategy`: (Optional) Strategy for managing conflicts. Can be "replace" or "merge" (default: "merge").
	// - `symlink`: (Optional) How to handle symlinks. Can be "deep", "shallow", or "skip" (default: "skip").
	//
	// ## Example
	//
	// ```lua
	//
	// 	copy({
	//  source = "./folder",
	//  dest = "./copy"
	// 	})
	// ```
	//
	Registry.RegisterModule(func(glue *core.Glue) error {
		glue.Plug("copy", core.MODULE).
			Brief("Copies folder").
			Arg("opts", CustomStruct("CopyOpts", []Field{
				NewField("source", STRING, "the file or folder to copy"),
				NewField("dest", STRING, "the destination to copy to"),
				NewField("strategy?", STRING, "a strategy for how to manage conflicts (replace or merge, defaults to merge)"),
				NewField("symlink?", STRING, "how to handle symlinks (deep/shallow/skip or the default skip)"),
			}), "the copy options").
			Do(func(R Runtime, args *Arguments) (RTValue, error) {
				opts, err := DecodeMap[CopyOpts](args.EnsureDict(0).Map())

				if err != nil {
					return nil, err
				}

				dest, err := glue.SmartPath(opts.Dest)

				if err != nil {
					return nil, err
				}

				src, err := glue.SmartPath(opts.Source)

				if err != nil {
					return nil, err
				}

				glue.Log.Info("[Copy]", "src", opts.Source, "dst", opts.Dest)

				opts.Dest = dest
				opts.Source = src

				return nil, Copy(opts)
			})

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
