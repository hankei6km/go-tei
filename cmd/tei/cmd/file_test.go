package cmd

import (
	"io/ioutil"
	"testing"

	"github.com/hankei6km/go-tei"
	"github.com/hankei6km/go-tei/cmd/tei/cli"
	"github.com/stretchr/testify/assert"
)

func Test_newFileCmd(t *testing.T) {
	teiBuilder := tei.NewBuilder()
	type args struct {
		teiBuilder tei.Builder
		cliBuilder cli.Builder
		args       []string
	}
	tests := []struct {
		name    string
		args    args
		want    cli.Cli
		wantErr bool
	}{
		{
			name: "basic",
			args: args{
				teiBuilder: teiBuilder,
				cliBuilder: NewFakeCliBuilder(),
				args:       []string{"/path/to/test.txt"},
			},
			want: cli.NewBuilder().
				CmdName("file").
				TeiBuilder(teiBuilder).
				OutStream(ioutil.Discard).
				ErrStream(ioutil.Discard).
				File("/path/to/test.txt").
				Build(),
		}, {
			name: "args=0",
			args: args{
				teiBuilder: teiBuilder,
				cliBuilder: NewFakeCliBuilder(),
				args:       []string{},
			},
			wantErr: true,
		}, {
			name: "args=2",
			args: args{
				teiBuilder: teiBuilder,
				cliBuilder: NewFakeCliBuilder(),
				args:       []string{"/path/to/test1.txt", "/path/to/test2.txt"},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			saveCmdExit := cmdExit
			defer func() {
				cmdExit = saveCmdExit
			}()
			cmdExit = func(int) {}
			var spy = func(b *fakeCliBuilder) {
				got := b.b.Build()
				assert.Equal(t, tt.want, got, "cli.Builder.Build() in cmd")
			}
			builders := func(tei.Builder, cli.Builder) (tei.Builder, cli.Builder) {
				return tt.args.teiBuilder, tt.args.cliBuilder.(FakeCliBuilder).SetSpy(spy)
			}
			c := newFileCmd(builders)
			c.SetArgs(tt.args.args)
			c.SetOutput(ioutil.Discard)
			err := c.Execute()
			if (err != nil) != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
