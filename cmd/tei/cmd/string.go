package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var newLineString = fmt.Sprintln()

func newStringCmd(builders globalBuildersFunc) *cobra.Command {
	var doNotNewLine bool
	// stringCmd represents the string command
	cmd := &cobra.Command{
		Use:          "string [flags] <string>...",
		SilenceUsage: true,
		Short:        "Switch the piped input to \"echo string\"",
		Long: `string switch the piped input to "echo string" if no data from the piped input.
`,
		Example: `  $ ` + cmdName + ` string "standby data"                         # standby data
  $ echo "" | ` + cmdName + ` string "standby data"               # standby data 
  $ echo "input data" | ` + cmdName + ` string "staandby data"    # input data`,
		Args: cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			teiBuilder, cliBuilder := builders(nil, nil)
			stringIntl := strings.Join(args, " ")
			if doNotNewLine == false {
				stringIntl = stringIntl + newLineString
			}
			stringCli := cliBuilder.
				CmdName(cmd.Name()).
				TeiBuilder(teiBuilder).
				String(stringIntl).
				InStream(cmd.InOrStdin()).
				OutStream(cmd.OutOrStdout()).
				ErrStream(cmd.ErrOrStderr()).
				Build()
			runCli(context.Background(), stringCli)
		},
	}

	flags := cmd.Flags()
	flags.SetInterspersed(false)

	flags.BoolVarP(&doNotNewLine, "n", "n", false, "do not output the trailing newline")
	return cmd
}

func init() {
	rootCmd.AddCommand(newStringCmd(builders))
}
