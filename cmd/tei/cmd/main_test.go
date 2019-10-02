package cmd

import (
	"io"
	"os"
	"runtime"
	"strings"
	"testing"

	"github.com/hankei6km/go-tei"
	"github.com/hankei6km/go-tei/cmd/tei/cli"
	"github.com/stretchr/testify/assert"
)

// main.go のテストの代わり.
// 実際にコマンドを実行した場合とは若干挙動が異なる(TODO を参照)

func Test_main(t *testing.T) {
	type args struct {
		args  []string
		input io.Reader
	}
	tests := []struct {
		name         string
		args         args
		wantOutText  string
		wantErrText  string
		wantExitCode int
	}{
		{
			name: "basic",
			args: args{
				args:  []string{"1"},
				input: os.Stdin,
			},
			wantExitCode: 0,
		}, {
			name: "some data",
			args: args{
				args:  []string{"1"},
				input: strings.NewReader("input data"),
			},
			wantExitCode: 1,
		}, {
			name: "no data",
			args: args{
				args:  []string{"1"},
				input: strings.NewReader(""),
			},
			wantExitCode: 0,
		}, {
			name: "leading newline CRLF",
			args: args{
				args:  []string{"1"},
				input: strings.NewReader("\r\n"),
			},
			wantExitCode: 0,
		}, {
			name: "leading newline LF",
			args: args{
				args:  []string{"1"},
				input: strings.NewReader("\n"),
			},
			wantExitCode: 0,
		}, {
			name: "leading newline CR",
			args: args{
				args:  []string{"1"},
				input: strings.NewReader("\r"),
			},
			wantExitCode: 0,
		}, {
			name: "leading newline CRLF + data",
			args: args{
				args:  []string{"1"},
				input: strings.NewReader("\r\nd"),
			},
			wantExitCode: 1,
		}, {
			name: "leading newline LF + data",
			args: args{
				args:  []string{"1"},
				input: strings.NewReader("\nd"),
			},
			wantExitCode: 1,
		}, {
			name: "leading newline CR + data",
			args: args{
				args:  []string{"1"},
				input: strings.NewReader("\rd"),
			},
			wantExitCode: 1,
		}, {
			name: "-l=false CRLF",
			args: args{
				args:  []string{"-l=false", "1"},
				input: strings.NewReader("\r\n"),
			},
			wantExitCode: 1,
		}, {
			name: "-l=false LF",
			args: args{
				args:  []string{"-l=false", "1"},
				input: strings.NewReader("\n"),
			},
			wantExitCode: 1,
		}, {
			name: "-l=false CR",
			args: args{
				args:  []string{"-l=false", "1"},
				input: strings.NewReader("\r"),
			},
			wantExitCode: 1,
		}, {
			name: "-p=true -l=false no data",
			args: args{
				args:  []string{"-p=true", "-l=false", "1"},
				input: strings.NewReader(""),
			},
			wantExitCode: 0,
		}, {
			name: "-p=true -l=false some data",
			args: args{
				args:  []string{"-p=true", "-l=false", "1"},
				input: strings.NewReader("input data"),
			},
			wantExitCode: 1,
			wantOutText:  "input data",
		}, {
			name: "-p=true -l=true some data",
			args: args{
				args:  []string{"-p=true", "-l=true", "1"},
				input: strings.NewReader("input data"),
			},
			wantExitCode: 2,
			wantErrText:  "Error: '-l' and '-p' flag has conflicted\n",
		}, {
			name: "parse err",
			args: args{
				args:  []string{"one"},
				input: strings.NewReader(""),
			},
			wantErrText: `Error in parsing exit code from args: strconv.Atoi: parsing "one": invalid syntax
`,
			wantExitCode: 128,
		}, {
			name: "file: basic",
			args: args{
				args:  []string{"file", testStandbyFile()},
				input: os.Stdin,
			},
			wantOutText: `file data
`,
		}, {
			name: "file: data",
			args: args{
				args:  []string{"file", testStandbyFile()},
				input: strings.NewReader("input data"),
			},
			wantOutText: `input data`,
		}, {
			name: "file: no data",
			args: args{
				args:  []string{"file", testStandbyFile()},
				input: strings.NewReader(""),
			},
			wantOutText: `file data
`,
		}, {
			name: "file: leading newline",
			args: args{
				args:  []string{"file", testStandbyFile()},
				input: strings.NewReader("\n"),
			},
			wantOutText: `file data
`,
		}, {
			name: "file: -l=false leading newline",
			args: args{
				args:  []string{"-l=false", "file", testStandbyFile()},
				input: strings.NewReader("\n"),
			},
			wantOutText: `
`,
		}, {
			name: "file: file not exist",
			args: args{
				args:  []string{"file", "not_exist"},
				input: os.Stdin,
			},
			wantErrText: `Error in runCli(file): Cli.Run reading the switched input: fileCli.Run open file: open ` +
				`not_exist: no such file or directory
`,
			wantExitCode: 1,
		}, {
			name: "file: args=0",
			args: args{
				args:  []string{"file"},
				input: os.Stdin,
			},
			// TODO: Args でのエラーが outSream になる原因調査(実際に動かすと stderr(2) へ出力されている).
			wantOutText: `Error: accepts 1 arg(s), received 0
`,
			wantExitCode: 1,
		}, {
			name: "string: basic",
			args: args{
				args:  []string{"string", "standby", "data"},
				input: os.Stdin,
			},
			wantOutText: `standby data
`,
		}, {
			name: "string: data",
			args: args{
				args:  []string{"string", "standby", "data"},
				input: strings.NewReader("input data"),
			},
			wantOutText: `input data`,
		}, {
			name: "string: no data",
			args: args{
				args:  []string{"string", "standby", "data"},
				input: strings.NewReader(""),
			},
			wantOutText: `standby data
`,
		}, {
			name: "string: leading newline",
			args: args{
				args:  []string{"string", "standby", "data"},
				input: strings.NewReader("\n"),
			},
			wantOutText: `standby data
`,
		}, {
			name: "string: l=false leading newline",
			args: args{
				args:  []string{"-l=false", "string", "standby", "data"},
				input: strings.NewReader("\n"),
			},
			wantOutText: `
`,
		}, {
			name: "string: -n",
			args: args{
				args:  []string{"string", "-n", "standby", "data"},
				input: strings.NewReader(""),
			},
			wantOutText: `standby data`,
		}, {
			name: "string: args=0",
			args: args{
				args:  []string{"string"},
				input: os.Stdin,
			},
			// TODO: Args でのエラーが outSream になる原因調査(実際に動かすと stderr(2) へ出力されている).
			wantOutText: `Error: requires at least 1 arg(s), only received 0
`,
			wantExitCode: 1,
		}, {
			name: "run: basic",
			args: args{
				args:  []string{"run", testStandbyCmd(), "test"},
				input: os.Stdin,
			},
			wantOutText: `standby cmd: test
`,
		}, {
			name: "run: data",
			args: args{
				args:  []string{"run", testStandbyCmd(), "test"},
				input: strings.NewReader("input data"),
			},
			wantOutText: `input data`,
		}, {
			name: "run: no data",
			args: args{
				args:  []string{"run", testStandbyCmd(), "test"},
				input: strings.NewReader(""),
			},
			wantOutText: `standby cmd: test
`,
		}, {
			name: "run: leading newline",
			args: args{
				args:  []string{"run", testStandbyCmd(), "test"},
				input: strings.NewReader("\n"),
			},
			wantOutText: `standby cmd: test
`,
		}, {
			name: "run: -l=false leading newline",
			args: args{
				args:  []string{"-l=false", "run", testStandbyCmd(), "test"},
				input: strings.NewReader("\n"),
			},
			wantOutText: `
`,
		}, {
			name: "run: error",
			args: args{
				args:  []string{"run", testStandbyCmdErr(), "test"},
				input: os.Stdin,
			},
			wantErrText: `Error in runCli(run): Cli.Run reading the switched input: runCli run - wait args([` +
				testStandbyCmdErr() +
				` test]): exit status 1
`,
			wantExitCode: 1,
		}, {
			name: "run: stderr",
			args: args{
				args:  []string{"run", testStandbyCmdErrOut(), "test"},
				input: os.Stdin,
			},
			wantErrText: `Error in runCli(run): Cli.Run reading the switched input: runCli run - errStream: standby cmd errout: test

`,
			wantExitCode: 1,
		}, {
			name: "run: command not exist",
			args: args{
				args:  []string{"run", "./not_exist"},
				input: os.Stdin,
			},
			wantErrText: `Error in runCli(run): Cli.Run reading the switched input: runCli run - start args([./not_exist]): fork/exec ./not_exist: no such file or directory
`,
			wantExitCode: 1,
		}, {
			name: "run: args=0",
			args: args{
				args:  []string{"run"},
				input: os.Stdin,
			},
			// TODO: Args でのエラーが outSream になる原因調査(実際に動かすと stderr(2) へ出力されている).
			wantOutText: `Error: requires at least 1 arg(s), only received 0
`,
			wantExitCode: 1,
		}, {
			name: "version",
			args: args{
				args:  []string{"version"},
				input: os.Stdin,
			},
			wantOutText: `tei
 version:    v99.99.99
 go version: ` + runtime.Version() + `
 os/arch:    ` + runtime.GOOS + "/" + runtime.GOARCH + `
 git commit: 1234abcd
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			saveCmdExit := cmdExit
			saveCmdVer := cmdVer
			saveCmdCommitHash := cmdCommitHash
			defer func() {
				cmdExit = saveCmdExit
				cmdVer = saveCmdVer
				cmdCommitHash = saveCmdCommitHash
			}()
			teiBuilder := tei.NewBuilder()
			cliBuilder := cli.NewBuilder()
			builders := func(t tei.Builder, c cli.Builder) (tei.Builder, cli.Builder) {
				if t != nil {
					teiBuilder = t
				}
				if c != nil {
					cliBuilder = c
				}
				return teiBuilder, cliBuilder
			}

			cmdVer = "v99.99.99"
			cmdCommitHash = "1234abcd"

			c := newRootCmd(builders)
			outStream := &strings.Builder{}
			errStream := &strings.Builder{}
			c.SetOut(outStream)
			c.SetErr(errStream)

			c.SetArgs(tt.args.args)
			c.SetIn(tt.args.input)

			c.AddCommand(newFileCmd(builders))
			c.AddCommand(newRunCmd(builders))
			c.AddCommand(newStringCmd(builders))
			c.AddCommand(newVersionCmd())

			cmdExit = func(exitCode int) {
				assert.Equal(t, tt.wantExitCode, exitCode, "exit code from newRootCmd().Execute()")
			}
			c.Execute()
			assert.Equal(t, tt.wantOutText, outStream.String(), "out from from newRootCmd().Execute()")
			assert.Equal(t, tt.wantErrText, errStream.String(), "err from from newRootCmd().Execute()")
		})
	}
}
