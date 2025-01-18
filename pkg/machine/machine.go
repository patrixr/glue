package machine

import (
	"io"
)

type Machine interface {
	Shell(input string, stdout io.Writer, stderr io.Writer) error
	TempFile(name string) (File, func() error, error)
}

type File interface {
	io.Writer
	io.StringWriter
	Name() string
}
