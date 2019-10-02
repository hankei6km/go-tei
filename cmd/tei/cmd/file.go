package cmd

import (
	"context"

	"github.com/spf13/cobra"
)

func newFileCmd(builders globalBuildersFunc) *cobra.Command {
	// fileCmd represents the file command
	cmd := &cobra.Command{
		Use:                   "file <input_file>",
		DisableFlagsInUseLine: true,
		SilenceUsage:          true,
		Short:                 "Switch the piped input to \"cat file\"",
		Long: `file switch the piped input to "cat /path/to/file" if no data from the piped input.
`,
		Example: `  $ ` + cmdName + ` file standby_data.txt                           # cat standby_data.txt
  $ echo -n "" | ` + cmdName + ` file standby_data.txt              # cat standby_data.txt
  $ echo -n "input data" | ` + cmdName + ` file standby_data.txt    # input data`,
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			teiBuilder, cliBuilder := builders(nil, nil)
			fileCli := cliBuilder.
				CmdName(cmd.Name()).
				TeiBuilder(teiBuilder).
				File(args[0]).
				InStream(cmd.InOrStdin()).
				OutStream(cmd.OutOrStdout()).
				ErrStream(cmd.ErrOrStderr()).
				Build()
			runCli(context.Background(), fileCli)
		},
	}

	flags := cmd.Flags()
	flags.SetInterspersed(false)

	return cmd
}

func init() {
	rootCmd.AddCommand(newFileCmd(builders))
}
