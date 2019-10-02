package cli

import (
	"context"
	"io"
	"os"

	"github.com/hankei6km/go-tei"
	"github.com/hankei6km/go-tei/internal/errors"
)

// Cli provides controler of command line tools.
//
// Cli はコマンドラインツールの制御を提供する.
type Cli interface {
	CmdName() string
	InStream() io.Reader
	OutStream() io.Writer
	ErrStream() io.Writer
	Run(ctx context.Context) (exitCode int, err error)
}

// Builder builds Cli.
type Builder interface {
	CmdName(string) Builder

	TeiBuilder(tei.Builder) Builder
	InStream(inStream io.Reader) Builder
	OutStream(outStream io.Writer) Builder
	ErrStream(errStream io.Writer) Builder

	FileCliBuilder
	RunCliBuilder
	StringCliBuilder

	Build() Cli
}

// builder impletems Builder.
//
// 今回は再利用等はしないので浅いコピーのみで対応.
type builder struct {
	cmdName string

	teiBuilder tei.Builder
	inStream   io.Reader
	outStream  io.Writer
	errStream  io.Writer

	file       string
	cmdArgs    []string
	stringIntl string
}

func (b *builder) CmdName(cmdName string) Builder {
	bb := b.branch()
	bb.cmdName = cmdName
	return bb
}

func (b *builder) TeiBuilder(teiBuilder tei.Builder) Builder {
	bb := b.branch()
	bb.teiBuilder = teiBuilder
	return bb
}

func (b *builder) InStream(inStream io.Reader) Builder {
	bb := b.branch()
	bb.inStream = inStream
	return bb
}

func (b *builder) OutStream(outStream io.Writer) Builder {
	bb := b.branch()
	bb.outStream = outStream
	return bb
}

func (b *builder) ErrStream(errStream io.Writer) Builder {
	bb := b.branch()
	bb.errStream = errStream
	return bb
}

func (b *builder) File(file string) Builder {
	bb := b.branch()
	bb.file = file
	return bb
}

func (b *builder) CmdArgs(cmdArgs []string) Builder {
	bb := b.branch()
	bb.cmdArgs = cmdArgs
	return bb
}

func (b *builder) String(stringIntl string) Builder {
	bb := b.branch()
	bb.stringIntl = stringIntl
	return bb
}

func (b *builder) branch() *builder {
	// return &(*b)
	return b // 今回は再利用の予定はないので、そのまま返す。
}
func (b *builder) Branch() Builder {
	return b.branch()
}

func (b *builder) Build() Cli {
	switch {
	case b.file != "":
		return newFileCli(b)
	case len(b.cmdArgs) > 0:
		return newRunCli(b)
	case b.stringIntl != "":
		return newStringCli(b)
	}
	return newBaseCli(b)
}

type baseCli struct {
	cmdName string

	teiBuilder tei.Builder
	inStream   io.Reader
	outStream  io.Writer
	errStream  io.Writer
}

func (c *baseCli) CmdName() string {
	return c.cmdName
}

func (c *baseCli) InStream() io.Reader {
	return c.inStream
}

func (c *baseCli) OutStream() io.Writer {
	return c.outStream
}

func (c *baseCli) Out() io.Writer {
	return c.outStream
}

func (c *baseCli) ErrStream() io.Writer {
	return c.errStream
}

func (c *baseCli) Run(ctx context.Context) (exitCode int, err error) {
	tei := c.teiBuilder.Build()
	_, err = io.Copy(c.outStream, tei.Switch(c.inStream))
	if err != nil {
		// TODO: err により exit code を変更.
		return 1, errors.Wrapf(err, "Cli.Run reading the switched input")
	}
	return
}

type errCli struct {
	baseCli
	exitCode int
	err      error
}

func (c *errCli) Run(ctx context.Context) (exitCode int, err error) {
	return c.exitCode, c.err
}

func newBaseCli(b *builder) *baseCli {
	return &baseCli{
		cmdName:    b.cmdName,
		teiBuilder: b.teiBuilder,
		inStream:   b.inStream,
		outStream:  b.outStream,
		errStream:  b.errStream,
	}
}

func newErrCli(b *builder, exitCode int, err error) *errCli {
	return &errCli{
		baseCli:  *newBaseCli(b),
		exitCode: exitCode,
		err:      err,
	}
}

func newCli(b *builder) Cli {
	return newBaseCli(b)
}

// NewBuilder returns the instance of Builder.
func NewBuilder() Builder {
	return &builder{
		teiBuilder: tei.NewBuilder(),
		inStream:   os.Stdin,
		outStream:  os.Stdout,
		errStream:  os.Stderr,
	}
}
