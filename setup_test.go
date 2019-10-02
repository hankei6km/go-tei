package tei

import (
	"os"
	"path/filepath"
	"testing"

	"go.uber.org/goleak"
)

func testDummyFile() *os.File {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	file, err := os.Open(filepath.Join(cwd, "testdata", "dummy.txt"))
	if err != nil {
		panic(err)
	}
	return file
}

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m)
}
