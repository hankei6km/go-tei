package cmd

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/hankei6km/go-tei"
	"github.com/hankei6km/go-tei/cmd/tei/cli"
	"github.com/hankei6km/go-tei/internal/errors"
	"go.uber.org/goleak"
)

type FakeCliBuilder interface {
	cli.Builder
	Act(bool) FakeCliBuilder
	ExitCode(int) FakeCliBuilder
	Err(error) FakeCliBuilder
	ErrText(string) FakeCliBuilder
	SetSpy(func(*fakeCliBuilder)) FakeCliBuilder // cli.Builder は浅いコピーなので必ずしも必要ではない.
}

type fakeCliBuilder struct {
	b        cli.Builder
	act      bool
	exitCode int
	err      error
	errText  string
	spy      func(*fakeCliBuilder)
}

func (b *fakeCliBuilder) CmdName(cmdName string) cli.Builder {
	b.b.CmdName(cmdName)
	return b
}

func (b *fakeCliBuilder) TeiBuilder(teiBuilder tei.Builder) cli.Builder {
	b.b.TeiBuilder(teiBuilder)
	return b
}

func (b *fakeCliBuilder) InStream(inStream io.Reader) cli.Builder {
	b.b.InStream(inStream)
	return b
}

func (b *fakeCliBuilder) OutStream(outStream io.Writer) cli.Builder {
	b.b.OutStream(outStream)
	return b
}

func (b *fakeCliBuilder) ErrStream(errStream io.Writer) cli.Builder {
	b.b.ErrStream(errStream)
	return b
}

func (b *fakeCliBuilder) File(file string) cli.Builder {
	b.b.File(file)
	return b
}

func (b *fakeCliBuilder) CmdArgs(cmdArgs []string) cli.Builder {
	b.b.CmdArgs(cmdArgs)
	return b
}

func (b *fakeCliBuilder) String(stringIntl string) cli.Builder {
	b.b.String(stringIntl)
	return b
}

func (b *fakeCliBuilder) Branch() cli.Builder {
	return b
}

func (b *fakeCliBuilder) Act(act bool) FakeCliBuilder {
	b.act = act
	return b
}

func (b *fakeCliBuilder) ExitCode(exitCode int) FakeCliBuilder {
	b.exitCode = exitCode
	return b
}

func (b *fakeCliBuilder) Err(err error) FakeCliBuilder {
	b.err = err
	return b
}

func (b *fakeCliBuilder) ErrText(errText string) FakeCliBuilder {
	b.errText = errText
	return b
}

func (b *fakeCliBuilder) SetSpy(spy func(*fakeCliBuilder)) FakeCliBuilder {
	b.spy = spy
	return b
}

func (b *fakeCliBuilder) Build() cli.Cli {
	b.spy(b)
	if b.act {
		return b.b.Build()
	}
	return NewFakeCli(b)
}

type fakeCli struct {
	b        cli.Cli
	exitCode int
	err      error
	errText  string
}

func (c *fakeCli) CmdName() string {
	return c.b.CmdName()
}

func (c *fakeCli) Args() []string {
	return []string{}
}

func (c *fakeCli) InStream() io.Reader {
	return c.b.InStream()
}

func (c *fakeCli) OutStream() io.Writer {
	return c.b.OutStream()
}

func (c *fakeCli) ErrStream() io.Writer {
	return c.b.ErrStream()
}

func (c *fakeCli) Run(ctx context.Context) (exitCode int, err error) {
	if _, err := io.Copy(c.b.ErrStream(), strings.NewReader(c.errText)); err != nil {
		return 1, errors.Wrapf(err, "fakeCli.Run copy errText to ErrStream")
	}
	return c.exitCode, c.err
}

func NewFakeCli(b *fakeCliBuilder) cli.Cli {
	return &fakeCli{
		b:        b.b.Build(),
		exitCode: b.exitCode,
		err:      b.err,
		errText:  b.errText,
	}
}

func NewFakeCliBuilder() FakeCliBuilder {
	return &fakeCliBuilder{
		b:   cli.NewBuilder(),
		spy: func(*fakeCliBuilder) {},
	}
}

// cli/setup_test.go と重複している.

func testDummyFile() string {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return filepath.Join(cwd, "testdata", "dummy.txt")
}

func testStandbyFile() string {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return filepath.Join(cwd, "testdata", "standby.txt")
}

func testStandbyCmd() string {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return filepath.Join(cwd, "testdata", "standby.sh")
}

func testStandbyCmdErr() string {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return filepath.Join(cwd, "testdata", "standby_err.sh")
}

func testStandbyCmdErrOut() string {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return filepath.Join(cwd, "testdata", "standby_err_out.sh")
}

func TestMain(m *testing.M) {
	saveCmdExit := cmdExit
	defer func() {
		cmdExit = saveCmdExit
	}()
	cmdExit = func(exitCode int) {
		panic("cmdExit: unexpected exit")
	}
	m.Run()
	goleak.VerifyTestMain(m)
}
