package modules

import (
	"errors"
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/patrixr/glue/pkg/core"
	. "github.com/patrixr/glue/pkg/runtime"
	"github.com/patrixr/q"
)

func init() {
	// @auteur("Modules/Blockinfile")
	//
	// # Blockinfile
	//
	// The blockinfile function allows you to insert, update, or remove a block of multi-line text in a file.
	// The block is surrounded by customizable markers to define its boundaries.
	//
	// Options:
	// - block: The multi-line text block to be inserted or updated.
	// - insertafter: A pattern to insert the block after.
	// - insertbefore: A pattern to insert the block before.
	// - marker: The marker template with "{mark}" as a placeholder.
	// - markerbegin: The text to replace "{mark}" for the beginning marker.
	// - markerend: The text to replace "{mark}" for the ending marker.
	// - state: A boolean to indicate whether to insert/update (true) or remove (false) the block.
	// - backup: A boolean to indicate whether to create a backup of the file.
	// - create: A boolean to indicate whether to create the file if it does not exist.
	//
	// Example usage:
	//
	// ```lua
	//
	// 	Blockinfile({
	// 		Insertafter = "pattern to insert after",
	// 		Insertbefore = "pattern to insert before",
	// 		Marker = "# {mark}",
	// 		Markerbegin = "BEGIN GLUE CUSTOM BLOCK",
	// 		Markerend = "END GLUE CUSTOM BLOCK",
	// 		State = true,
	// 		Backup = true,
	// 		Create = true,
	// 		Block = [[
	// 			This is a block of text
	// 			that spans multiple lines.
	// 		]],
	// 	})
	// ```
	Registry.RegisterModule(func(glue *core.Glue) error {
		glue.Plug("blockinfile", core.MODULE).
			Brief("Insert/update/remove a block of multi-line text surrounded by customizable markers in a file").
			Arg("block_params", CustomStruct("BlockinfileParams", []Field{
				NewField("path", STRING, "the file to insert the block into"),
				NewField("block", STRING, "the multi-line text block to be inserted or updated"),
				NewField("insertafter?", STRING, "the multi-line text block to be inserted or updated"),
				NewField("insertbefore?", STRING, "the multi-line text block to be inserted or updated"),
				NewField("marker?", STRING, "the multi-line text block to be inserted or updated"),
				NewField("markerbegin?", STRING, "the multi-line text block to be inserted or updated"),
				NewField("markerend?", STRING, "the multi-line text block to be inserted or updated"),
				NewField("state", BOOL, "the multi-line text block to be inserted or updated"),
				NewField("backup?", BOOL, "the multi-line text block to be inserted or updated"),
				NewField("create?", BOOL, "the multi-line text block to be inserted or updated"),
			}), "the configuration for the block insertion").
			Do(func(R Runtime, args *Arguments) (RTValue, error) {
				data := args.EnsureDict(0).Map()
				props, err := DecodeMap[BlockOpts](data)

				if err != nil {
					return nil, err
				}

				path, err := glue.SmartPath(props.Path)

				if err != nil {
					return nil, err
				}

				glue.Log.Info("[Blockinfile]", "path", props.Path)

				props.Path = path

				return nil, BlockInFile(props)
			})

		return nil
	})
}

const defaultMarker = "# {mark}"
const defaultMarkerBegin = "BEGIN GLUE MANAGED BLOCK"
const defaultMarkerEnd = "END GLUE MANAGED BLOCK"

type BlockOpts struct {
	Block        string `json:"block"`
	Insertafter  string `json:"insertafter"`
	Insertbefore string `json:"insertbefore"`
	Marker       string `json:"marker"`
	Markerbegin  string `json:"markerbegin"`
	Markerend    string `json:"markerend"`
	State        bool   `json:"state"`
	Backup       bool   `json:"backup"`
	Create       bool   `json:"create"`
	Path         string `json:"path"`
}

func BlockInString(text string, opts BlockOpts) string {
	eol := "\n"
	markerBegin := stringOr(opts.Markerbegin, defaultMarkerBegin)
	markerEnd := stringOr(opts.Markerend, defaultMarkerEnd)
	marker := stringOr(opts.Marker, defaultMarker)
	beginLine := strings.Replace(marker, "{mark}", markerBegin, 1)
	endLine := strings.Replace(marker, "{mark}", markerEnd, 1)
	lines := strings.Split(text, eol)
	found, _, beginIdx := q.Find(lines, q.Eq(beginLine))
	foundEnd, _, endIdx := q.Find(lines, q.Eq(endLine))
	blockExists := found && foundEnd
	block := opts.Block

	if !opts.State {
		if !blockExists {
			// block is already absent from text
			return text
		}

		// remove existing block
		return strings.Join(append(lines[:beginIdx], lines[endIdx+1:]...), eol)
	}

	if found {
		// replace existing block
		before := strings.Join(lines[:beginIdx+1], eol)
		after := strings.Join(lines[endIdx:], eol)
		return before + eol + block + eol + after
	}

	insertAt := -1

	// prioritize insertafter if present
	if len(opts.Insertafter) > 0 {
		if found, _, idx := q.Find(lines, q.Match(opts.Insertafter)); found {
			insertAt = idx + 1
		}
	}

	// look for insertbefore if we don't already have a location
	if insertAt < 0 && len(opts.Insertbefore) > 0 {
		if found, _, idx := q.Find(lines, q.Match(opts.Insertbefore)); found {
			insertAt = idx
		}
	}

	// default: place at the end
	if insertAt < 0 {
		insertAt = len(lines)
	}

	var tail string
	var head string = strings.Join(lines[:insertAt], eol)

	if insertAt < len(lines) {
		tail = strings.Join(lines[insertAt:], eol)
	}

	result := head + eol + strings.Join([]string{
		beginLine,
		block,
		endLine,
	},
		eol,
	) + eol + tail

	return regexp.MustCompile("\n*$").ReplaceAllString(result, eol)
}

func BlockInFile(props BlockOpts) error {
	path := props.Path
	stat, err := os.Stat(path)
	mode := os.FileMode(0644)
	flag := os.O_RDWR

	if props.State && len(props.Block) == 0 {
		return errors.New("Cannot insert empty block")
	}

	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}

		if !props.Create {
			return err
		}

		flag = flag | os.O_CREATE
	} else {
		mode = stat.Mode()
	}

	file, err := os.OpenFile(path, flag, mode)

	if err != nil {
		return err
	}

	defer file.Close()

	data, err := io.ReadAll(file)

	if err != nil {
		return err
	}

	source := string(data)
	updated := BlockInString(source, props)

	if props.Backup && updated != source {
		Backup(props.Path)
	}

	file.Truncate(0)
	file.Seek(0, 0)

	_, err = file.Write([]byte(updated))

	if err != nil {
		return err
	}

	return nil
}

func stringOr(txt string, fallback string) string {
	if len(txt) > 0 {
		return txt
	}
	return fallback
}
