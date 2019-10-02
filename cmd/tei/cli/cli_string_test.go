package cli

import (
	"context"
	"errors"
	"os"
	"strings"
	"testing"

	"github.com/hankei6km/go-tei"
	"github.com/stretchr/testify/assert"
)

func Test_stringCli_Run(t *testing.T) {
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
				String("string data"),
			args: args{
				ctx: context.Background(),
			},
			want: "input data",
		}, {
			name: "staandby",
			builder: NewBuilder().
				InStream(strings.NewReader("")).
				String("string data"),
			args: args{
				ctx: context.Background(),
			},
			want: "string data",
		}, {
			name: "stdin",
			builder: NewBuilder().
				InStream(os.Stdin).
				String("string data"),
			args: args{
				ctx: context.Background(),
			},
			want: "string data",
		}, {
			name: "error",
			builder: NewBuilder().
				InStream(tei.ErrReader(errors.New("test error"))).
				String("string data"),
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
				t.Errorf("stringCli.Run() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotExitCode != tt.wantExitCode {
				t.Errorf("stringCli.Run() = %v, want %v", gotExitCode, tt.wantExitCode)
			}
			assert.Equal(t, tt.want, outStream.String(), "stringCli.Run() outStream")
		})
	}
}
