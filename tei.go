// Package tei is an io.Reader switcher.
package tei

import (
	"bytes"
	"io"
	"os"

	"github.com/hankei6km/go-tei/internal/errors"
)

// StandbyFunc returns io.Reader built by the standby source.
type StandbyFunc func() io.Reader

// ErrReader returns io.Reader that retuns the err instead of io.EOF.
func ErrReader(err error) io.Reader {
	// TODO: err 伝播用の reader を作る?
	r, w := io.Pipe()
	go func() {
		w.CloseWithError(err)
	}()
	return r
}

// Tei is an io.Reader switcher.
type Tei interface {
	// Switch activates the standby source instead of the input, if nothing data from the input.
	Switch(input io.Reader) (r io.Reader)
}

// Builder builds Tei.
type Builder interface {
	// Standy sets the function that return the standby source.
	Standby(StandbyFunc) Builder
	// IgnoreLeadingNewline sets the flag that ignore leading a newline while sniffing the input.
	IgnoreLeadingNewline(bool) Builder
	// SwitchByTerminal sets the flag that force switch if the input is opened on a terminal.
	SwitchByTerminal(bool) Builder

	// Branch() Builder

	// Build builds the instance of Tei.
	Build() Tei
}

type baseBuilder struct {
	standby              StandbyFunc
	ignoreLeadingNewline bool
	switchByTerminal     bool
}

func (b *baseBuilder) Standby(standby StandbyFunc) Builder {
	bb := b.branch()
	bb.standby = standby
	return bb
}

func (b *baseBuilder) IgnoreLeadingNewline(ignoreLeadingNewline bool) Builder {
	bb := b.branch()
	bb.ignoreLeadingNewline = ignoreLeadingNewline
	return bb
}

func (b *baseBuilder) SwitchByTerminal(switchByTerminal bool) Builder {
	bb := b.branch()
	bb.switchByTerminal = switchByTerminal
	return bb
}

func (b *baseBuilder) branch() *baseBuilder {
	// return &(*b)
	return b // 今回は再利用の予定はないので、そのまま返す。
}

func (b *baseBuilder) Branch() Builder {
	return b.branch()
}

func (b *baseBuilder) Build() Tei {
	return newBaseTei(b)
}

type baseTei struct {
	standby              StandbyFunc
	ignoreLeadingNewline bool
	switchByTerminal     bool
}

func (t *baseTei) Switch(input io.Reader) (r io.Reader) {
	if t.switchByTerminal {
		if file, ok := input.(*os.File); ok {
			stat, err := file.Stat()
			if err != nil {
				return ErrReader(errors.Wrapf(err, "baseTei.Switch checking switchByTerminal"))
			}
			if stat.Mode()&os.ModeDevice != 0 {
				return t.standby()
			}
		}
	}

	buf := bytes.NewBuffer([]byte{})
	n, err := io.CopyN(buf, input, 3)
	switch {
	case err == io.EOF:
		switch {
		case n == 0:
			return t.standby()
		case t.ignoreLeadingNewline:
			b := buf.Bytes()
			switch {
			case n == 2 && b[0] == '\r' && b[1] == '\n':
				return t.standby()
			case n == 1 && b[0] == '\n':
				return t.standby()
			case n == 1 && b[0] == '\r':
				return t.standby()
			}
		}
	case err != nil && err != io.EOF:
		return ErrReader(errors.Wrapf(err, "baseTei.Switch sniffing the inpu"))
	}
	return io.MultiReader(buf, input)
}

func newBaseTei(b *baseBuilder) *baseTei {
	return &baseTei{
		standby:              b.standby,
		ignoreLeadingNewline: b.ignoreLeadingNewline,
		switchByTerminal:     b.switchByTerminal,
	}
}

// NewBuilder returns the instance of Builder.
func NewBuilder() Builder {
	return &baseBuilder{
		ignoreLeadingNewline: true,
		switchByTerminal:     true,
	}
}
