package cli

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_builder_Build(t *testing.T) {
	tests := []struct {
		name    string
		builder Builder
		want    Cli
	}{
		{
			name:    "bssic",
			builder: NewBuilder(),
			want:    &baseCli{},
		}, {
			name:    "file",
			builder: NewBuilder().File("test"),
			want:    &fileCli{},
		}, {
			name:    "run",
			builder: NewBuilder().CmdArgs([]string{"test"}),
			want:    &runCli{},
		}, {
			name:    "string",
			builder: NewBuilder().String("test"),
			want:    &stringCli{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.builder.Build()
			assert.IsType(t, tt.want, got, "builder.Build()")
		})
	}
}
