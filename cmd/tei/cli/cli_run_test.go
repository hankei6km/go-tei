package cli

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_runCli_Run(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name         string
		builder      Builder
		args         args
		want         string
		wantExitCode int
		wantErr      bool
	}{
		{
			name: "basic",
			builder: NewBuilder().
				InStream(strings.NewReader("input data")).
				CmdArgs([]string{testStandbyCmd()}),
			args: args{
				ctx: context.Background(),
			},
			want: "input data",
		}, {
			name: "stanby",
			builder: NewBuilder().
				InStream(strings.NewReader("")).
				CmdArgs([]string{testStandbyCmd(), "test"}),
			args: args{
				ctx: context.Background(),
			},
			want: "standby cmd: test\n",
		}, {
			name: "stdin",
			builder: NewBuilder().
				InStream(os.Stdin).
				CmdArgs([]string{testStandbyCmd(), "test"}),
			args: args{
				ctx: context.Background(),
			},
			want: "standby cmd: test\n",
		}, {
			name: "error",
			builder: NewBuilder().
				InStream(strings.NewReader("")).
				CmdArgs([]string{testStandbyCmdErr(), "test"}),
			args: args{
				ctx: context.Background(),
			},
			want:         "",
			wantExitCode: 1,
			wantErr:      true,
		}, {
			name: "error out",
			builder: NewBuilder().
				InStream(strings.NewReader("")).
				CmdArgs([]string{testStandbyCmdErrOut(), "test"}),
			args: args{
				ctx: context.Background(),
			},
			want:         "",
			wantExitCode: 1,
			wantErr:      true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			outStream := &strings.Builder{}
			c := tt.builder.
				OutStream(outStream).
				Build()
			gotExitCode, err := c.Run(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("runCli.Run() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotExitCode != tt.wantExitCode {
				t.Errorf("runCli.Run() = %v, want %v", gotExitCode, tt.wantExitCode)
			}
			assert.Equal(t, tt.want, outStream.String(), "runCli.Run() outStream")
		})
	}
}
