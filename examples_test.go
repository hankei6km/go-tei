package tei_test

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/hankei6km/go-tei"
)

func ExampleTei() {
	output := func(builder tei.Builder, r io.Reader) {
		tei := builder.Standby(func() io.Reader {
			return strings.NewReader("standby-data")
		}).Build()

		r = tei.Switch(r)
		io.Copy(os.Stdout, r)
		fmt.Println()
	}
	output(tei.NewBuilder(), os.Stdin)
	output(tei.NewBuilder(), strings.NewReader(""))
	output(tei.NewBuilder(), strings.NewReader("\n"))
	output(tei.NewBuilder(), strings.NewReader("input-data"))

	// Output:
	//
	// standby-data
	// standby-data
	// standby-data
	// input-data
}

func ExampleBuilder_IgnoreLeadingNewline() {
	output := func(builder tei.Builder, r io.Reader) {
		tei := builder.Standby(func() io.Reader {
			return nil
		}).Build()

		if tei.Switch(r) == nil {
			fmt.Println("no-data")
			return
		}
		fmt.Println("data")

	}
	output(tei.NewBuilder(), strings.NewReader("\n"))
	output(tei.NewBuilder(), strings.NewReader("\ninput-data"))
	output(tei.NewBuilder().IgnoreLeadingNewline(false), strings.NewReader("\n"))

	// Output:
	//
	// no-data
	// data
	// data
}
