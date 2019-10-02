package cmd

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/hankei6km/go-tei/cmd/tei/cli"
	"github.com/stretchr/testify/assert"
)

func Test_runCli(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name         string
		cliBuilder   cli.Builder
		args         args
		wantExitCode int
		wantErrOut   string
	}{
		{
			name:       "basic",
			cliBuilder: NewFakeCliBuilder(),
			args: args{
				ctx: context.Background(),
			},
		}, {
			name: "error",
			cliBuilder: NewFakeCliBuilder().CmdName("errCmd").(FakeCliBuilder).
				ExitCode(99).Err(errors.New("test err")),
			args: args{
				ctx: context.Background(),
			},
			wantExitCode: 99,
			wantErrOut:   "Error in runCli(errCmd): test err\n",
		}, {
			name: "err out",
			cliBuilder: NewFakeCliBuilder().CmdName("errCmd").(FakeCliBuilder).
				ExitCode(99).Err(errors.New("test err")).ErrText("test err out\n"),
			args: args{
				ctx: context.Background(),
			},
			wantExitCode: 99,
			wantErrOut:   "test err out\nError in runCli(errCmd): test err\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			saveExitBare := cmdExit
			defer func() { cmdExit = saveExitBare }()

			gotErrOut := &strings.Builder{}
			cliBuilder := tt.cliBuilder.ErrStream(gotErrOut)
			cmdExit = func(exitCode int) {
				assert.Equal(t, tt.wantExitCode, exitCode, "runCli() exitcode")
			}
			runCli(tt.args.ctx, cliBuilder.Build())
			assert.Equal(t, tt.wantErrOut, gotErrOut.String(), "runCli() errOut")
		})
	}
}
