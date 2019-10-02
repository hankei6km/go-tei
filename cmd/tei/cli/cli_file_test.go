package cli

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_fileCli_Run(t *testing.T) {
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
				File(testStandbyFile()),
			args: args{
				ctx: context.Background(),
			},
			want: "input data",
		}, {
			name: "staandby",
			builder: NewBuilder().
				InStream(strings.NewReader("")).
				File(testStandbyFile()),
			args: args{
				ctx: context.Background(),
			},
			want: "file data\n",
		}, {
			name: "stdin",
			builder: NewBuilder().
				InStream(os.Stdin).
				File(testStandbyFile()),
			args: args{
				ctx: context.Background(),
			},
			want: "file data\n",
		}, {
			name: "not exist",
			builder: NewBuilder().
				InStream(strings.NewReader("")).
				File(filepath.Join(testDummyFile(), "not_exist")),
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
				t.Errorf("fileCli.Run() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotExitCode != tt.wantExitCode {
				t.Errorf("fileCli.Run() = %v, want %v", gotExitCode, tt.wantExitCode)
			}
			assert.Equal(t, tt.want, outStream.String(), "fileCli.Run() outStream")
		})
	}
}
