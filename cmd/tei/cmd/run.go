package cmd

import (
	"context"

	"github.com/spf13/cobra"
)

func newRunCmd(builders globalBuildersFunc) *cobra.Command {
	// runCmd represents the run command
	cmd := &cobra.Command{
		Use:                   "run <command> [command_args]...",
		DisableFlagsInUseLine: true,
		SilenceUsage:          true,
		Short:                 "Switch the piped input to \"run command\"",
		Long: `run switch the piped input to "run command [command_args]..."
if no data from the piped input.
`,
		Example: `  $ ` + cmdName + ` run echo "standby data"                        # standby data
  $ echo "" | ` + cmdName + ` run echo "standby data"              # standby data
  $ echo "input data" | ` + cmdName + ` run echo "standby data"    # input data`,
		Args: cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			teiBuilder, cliBuilder := builders(nil, nil)
			rCli := cliBuilder.
				CmdName(cmd.Name()).
				TeiBuilder(teiBuilder).
				CmdArgs(args).
				InStream(cmd.InOrStdin()).
				OutStream(cmd.OutOrStdout()).
				ErrStream(cmd.ErrOrStderr()).
				Build()
			runCli(context.Background(), rCli)
		},
	}

	flags := cmd.Flags()
	flags.SetInterspersed(false)

	return cmd
}

func init() {
	rootCmd.AddCommand(newRunCmd(builders))
}
