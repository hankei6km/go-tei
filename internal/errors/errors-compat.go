// +build !go1.13

package errors

import (
	"github.com/pkg/errors"
)

// Wrapf wraps github.com/pkg/errors/Wrapf.
func Wrapf(err error, format string, a ...interface{}) error {
	return errors.Wrapf(err, format, a...)
}
