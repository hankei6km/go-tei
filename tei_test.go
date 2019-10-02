package tei

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_errReader(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name string
		args args
		want error
	}{
		{
			name: "basic",
			args: args{
				err: errors.New("test error"),
			},
			want: errors.New("test error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := io.Copy(ioutil.Discard, ErrReader(tt.args.err))
			assert.Equal(t, tt.want, err, "errReader() read error")
		})
	}
}

func Test_baseTei_Reader(t *testing.T) {
	standbyFunc := func() io.Reader {
		return bytes.NewBuffer([]byte("standby data"))
	}
	tests := []struct {
		name    string
		builder Builder
		input   io.Reader
		want    []byte
		wantErr bool
	}{
		{
			name: "basic",
			builder: NewBuilder().
				Standby(standbyFunc),
			input: bytes.NewBuffer([]byte("input data")),
			want:  []byte("input data"),
		}, {
			name: "standby",
			builder: NewBuilder().
				Standby(standbyFunc),
			input: bytes.NewBuffer([]byte{}),
			want:  []byte("standby data"),
		}, {
			name: "stdin",
			builder: NewBuilder().
				Standby(standbyFunc),
			input: os.Stdin,
			want:  []byte("standby data"),
		}, {
			name: "LF",
			builder: NewBuilder().
				Standby(standbyFunc),
			input: bytes.NewBuffer([]byte("\n")),
			want:  []byte("standby data"),
		}, {
			name: "CR",
			builder: NewBuilder().
				Standby(standbyFunc),
			input: bytes.NewBuffer([]byte("\r")),
			want:  []byte("standby data"),
		}, {
			name: "CRLF",
			builder: NewBuilder().
				Standby(standbyFunc),
			input: bytes.NewBuffer([]byte("\r\n")),
			want:  []byte("standby data"),
		}, {
			name: "CR+data",
			builder: NewBuilder().
				Standby(standbyFunc),
			input: bytes.NewBuffer([]byte("\rT")),
			want:  []byte("\rT"),
		}, {
			name: "CRLF+data",
			builder: NewBuilder().
				Standby(standbyFunc),
			input: bytes.NewBuffer([]byte("\r\nT")),
			want:  []byte("\r\nT"),
		}, {
			name: "ignoreLeadingNewline=false LF",
			builder: NewBuilder().
				Standby(standbyFunc).
				IgnoreLeadingNewline(false),
			input: bytes.NewBuffer([]byte("\n")),
			want:  []byte("\n"),
		}, {
			name: "ignoreLeadingNewline=false CR",
			builder: NewBuilder().
				Standby(standbyFunc).
				IgnoreLeadingNewline(false),
			input: bytes.NewBuffer([]byte("\r")),
			want:  []byte("\r"),
		}, {
			name: "ignoreLeadingNewline=false CRLF",
			builder: NewBuilder().
				Standby(standbyFunc).
				IgnoreLeadingNewline(false),
			input: bytes.NewBuffer([]byte("\r\n")),
			want:  []byte("\r\n"),
		}, {
			name: "ignoreLeadingNewline=false CR+data",
			builder: NewBuilder().
				Standby(standbyFunc).
				IgnoreLeadingNewline(false),
			input: bytes.NewBuffer([]byte("\rT")),
			want:  []byte("\rT"),
		}, {
			name: "ignoreLeadingNewline=false CRLF+data",
			builder: NewBuilder().
				Standby(standbyFunc).
				IgnoreLeadingNewline(false),
			input: bytes.NewBuffer([]byte("\r\nT")),
			want:  []byte("\r\nT"),
		}, {
			name: "switchByTerminal=false",
			builder: NewBuilder().
				Standby(standbyFunc).
				SwitchByTerminal(false),
			input: os.Stdin,
			want:  []byte("standby data"),
		}, {
			name: "null",
			builder: NewBuilder().
				Standby(standbyFunc),
			input: bytes.NewBuffer([]byte{0}),
			want:  []byte{0},
		}, {
			name: "error",
			builder: NewBuilder().
				Standby(standbyFunc),
			input:   ErrReader(errors.New("test error")),
			want:    []byte{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := tt.builder.Build().Switch(tt.input)
			buf := bytes.NewBuffer([]byte{})
			_, err := io.Copy(buf, r)
			got := buf.Bytes()
			assert.Equal(t, tt.want, got, "baseTei.Switch() read", string(got))
			if assert.Equal(t, tt.wantErr, (err != nil), "baseTei.Reder() error in read") == false {
				log.Println(err)
			}
		})
	}
}
