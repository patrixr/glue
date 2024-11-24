package homebrew

import (
	"fmt"
	"io"
	"os"

	"time"

	"github.com/patrixr/glue/pkg/shell"
)

type Row struct {
	kind string
	name string
}

type Homebrew struct {
	rows   []Row
	stdout io.Writer
	stderr io.Writer
}

func NewHomebrew(stdout io.Writer, stderr io.Writer) *Homebrew {
	return &Homebrew{
		stdout: stdout,
		stderr: stderr,
	}
}

func (h *Homebrew) Package(pkg string) {
	h.rows = append(h.rows, Row{kind: "brew", name: pkg})
}

func (h *Homebrew) Cask(pkg string) {
	h.rows = append(h.rows, Row{kind: "cask", name: pkg})
}

func (h *Homebrew) Tap(pkg string) {
	h.rows = append(h.rows, Row{kind: "tap", name: pkg})
}

func (h *Homebrew) Mas(pkg string) {
	h.rows = append(h.rows, Row{kind: "mas", name: pkg})
}

func (h *Homebrew) Whalebrew(pkg string) {
	h.rows = append(h.rows, Row{kind: "whalebrew", name: pkg})
}

func (h *Homebrew) Install() error {
	tmp, err := os.CreateTemp("", ".glue_brewfile_"+time.Now().Format("20060102150405"))

	defer os.Remove(tmp.Name())

	if err != nil {
		return err
	}

	for _, row := range h.rows {
		tmp.WriteString(fmt.Sprintf("%s \"%s\"\n", row.kind, row.name))
	}

	return shell.Run(fmt.Sprintf("brew bundle --file=%s --no-lock", tmp.Name()), h.stdout, h.stderr)
}

func (h *Homebrew) Upgrade() error {
	return shell.Run("brew upgrade", h.stdout, h.stderr)
}
