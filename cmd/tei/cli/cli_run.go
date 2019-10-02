package cli

import (
	"context"
	"fmt"
	"io"
	"os/exec"
	"strings"

	"github.com/hankei6km/go-tei/internal/errors"
)

type runCli struct {
	baseCli
	cmdArgs []string
}

// RunCliBuilder adds a property to CliBuilder.
type RunCliBuilder interface {
	CmdArgs([]string) Builder
}

func (c *runCli) Run(ctx context.Context) (exitCode int, err error) {
	c.teiBuilder = c.teiBuilder.Standby(func() io.Reader {
		r, w := io.Pipe()
		go func(w *io.PipeWriter) {
			var cmdErr error
			errStream := &strings.Builder{}
			defer func() {
				switch {
				case cmdErr != nil:
					w.CloseWithError(cmdErr)
					return
				case errStream.Len() > 0:
					w.CloseWithError(errors.Wrapf(fmt.Errorf(errStream.String()), "runCli run - errStream"))
					return
				}
				w.Close()
			}()
			cmdPath := c.cmdArgs[0]
			var cmdArgs []string
			if len(c.cmdArgs) > 1 {
				cmdArgs = c.cmdArgs[1:]
			} else {
				cmdArgs = []string{}
			}
			cmd := exec.CommandContext(ctx, cmdPath, cmdArgs...)
			cmd.Stdout = w
			cmd.Stderr = errStream
			if err := cmd.Start(); err != nil {
				cmdErr = errors.Wrapf(err, "runCli run - start args(%s)", c.cmdArgs)
				return
			}
			if err := cmd.Wait(); err != nil {
				cmdErr = errors.Wrapf(err, "runCli run - wait args(%s)", c.cmdArgs)
			}

		}(w)
		return r
	})
	return c.baseCli.Run(ctx)
}

func newRunCli(b *builder) *runCli {
	return &runCli{
		baseCli: *newBaseCli(b),
		cmdArgs: b.cmdArgs,
	}
}
