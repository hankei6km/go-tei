package cli

import (
	"context"
	"io"
	"os"

	"github.com/hankei6km/go-tei"
	"github.com/hankei6km/go-tei/internal/errors"
)

type fileCli struct {
	baseCli
	file string
}

// FileCliBuilder adds a property to CliBuilder.
type FileCliBuilder interface {
	File(string) Builder
}

func (c *fileCli) Run(ctx context.Context) (exitCode int, err error) {
	var file *os.File
	defer func() {
		if file != nil {
			file.Close()
		}
	}()
	c.teiBuilder = c.teiBuilder.Standby(func() io.Reader {
		file, err = os.Open(c.file)
		if err != nil {
			file = nil
			return tei.ErrReader(errors.Wrapf(err, "fileCli.Run open file"))
		}
		return file
	})
	return c.baseCli.Run(ctx)
}

func newFileCli(b *builder) *fileCli {
	return &fileCli{
		baseCli: *newBaseCli(b),
		file:    b.file,
	}
}
