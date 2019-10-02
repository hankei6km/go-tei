package cli

import (
	"os"
	"path/filepath"
	"testing"

	"go.uber.org/goleak"
)

// cmd/setup_test.go と重複している.

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
	goleak.VerifyTestMain(m)
}
