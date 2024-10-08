package homebrew

import (
	"fmt"
	"os"

	"time"

	"github.com/patrixr/glue/pkg/shell"
)

type Row struct {
	kind string
	name string
}

type Homebrew struct {
	rows []Row
}

func NewHomebrew() *Homebrew {
	return &Homebrew{}
}

func (h *Homebrew) Brew(pkg string) {
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

	return shell.Run(fmt.Sprintf("brew bundle --file=%s --no-lock", tmp.Name()), os.Stdout, os.Stderr)
}
