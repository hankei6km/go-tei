package cli

import (
	"context"
	"io"
	"strings"
)

type stringCli struct {
	baseCli
	stringIntl string
}

// StringCliBuilder adds a property to CliBuilder.
type StringCliBuilder interface {
	String(string) Builder
}

func (c *stringCli) Run(ctx context.Context) (exitCode int, err error) {
	c.teiBuilder = c.teiBuilder.Standby(func() io.Reader {
		return strings.NewReader(c.stringIntl)
	})
	return c.baseCli.Run(ctx)
}

func newStringCli(b *builder) *stringCli {
	return &stringCli{
		baseCli:    *newBaseCli(b),
		stringIntl: b.stringIntl,
	}
}
