package cmd

import (
	"context"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/hankei6km/go-tei"
	"github.com/hankei6km/go-tei/cmd/tei/cli"
	"github.com/spf13/cobra"
)

var cmdName = "tei"
var cmdVer string
var cmdCommitHash string

type exitFunc func(exitCode int)

var cmdExit = func(exitCode int) {
	os.Exit(exitCode)
}

func runCli(ctx context.Context, cli cli.Cli) {
	exitCode, err := cli.Run(ctx)
	if err != nil {
		fmt.Fprintf(cli.ErrStream(), "Error in runCli(%s): %s\n", cli.CmdName(), err.Error())
	}
	cmdExit(exitCode)
}

type globalBuildersFunc func(tei.Builder, cli.Builder) (tei.Builder, cli.Builder)

var builders = func() globalBuildersFunc {
	var teiBuilder = tei.NewBuilder()
	var cliBuilder = cli.NewBuilder()
	return func(t tei.Builder, c cli.Builder) (tei.Builder, cli.Builder) {
		if t != nil {
			teiBuilder = t
		}
		if c != nil {
			cliBuilder = c
		}
		return teiBuilder, cliBuilder
	}
}()

func newRootCmd(builders globalBuildersFunc) *cobra.Command {
	var ignoreNewline bool
	var passThrough bool
	// rootCmd represents the base command when called without any subcommands
	cmd := &cobra.Command{
		Use:   "tei [flags] <exit_code>",
		Short: cmdName + " switch the piped input to another one if no data from the piped input",
		Long: cmdName + ` switch the piped input to another one if no data from the piped input,
and simply use to just check no data from the piped input.
`,
		Example: `  $ ` + cmdName + ` 1                        # exit code = 0
  $ echo "" | ` + cmdName + ` 1              # exit code = 0
  $ echo "input data" | ` + cmdName + ` 1    # exit code = 1`,
		Args: cobra.ExactArgs(1),
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			teiBuilder, _ := builders(nil, nil)
			teiBuilder = teiBuilder.IgnoreLeadingNewline(ignoreNewline)
			builders(teiBuilder, nil)
		},
		Run: func(cmd *cobra.Command, args []string) {
			exitCode, err := strconv.Atoi(args[0])
			if err != nil {
				fmt.Fprintf(cmd.ErrOrStderr(), "Error in parsing exit code from args: %s\n", err.Error())
				cmdExit(128) // 何で終了するのがよい?: https://www.tldp.org/LDP/abs/html/exitcodes.html
				return
			}
			if passThrough && ignoreNewline {
				fmt.Fprintf(cmd.ErrOrStderr(), "Error: '-l' and '-p' flag has conflicted\n")
				// cmdExit(1) exitCode = 1 と見分けがつかない
				cmdExit(exitCode + 1) // とりあえず.
				return
			}
			teiBuilder, _ := builders(nil, nil)
			tei := teiBuilder.Standby(func() io.Reader { return nil }).Build()
			inStream := tei.Switch(cmd.InOrStdin())
			if inStream != nil {
				if passThrough {
					if _, err := io.Copy(cmd.OutOrStdout(), inStream); err != nil {
						fmt.Fprintf(cmd.ErrOrStderr(), "Error in coping the input to stdout: %s\n", err.Error())
						// cmdExit(1) exitCode = 1 と見分けがつかない
						cmdExit(exitCode + 1) // とりあえず.
						return
					}
				}
				cmdExit(exitCode)
				return
			}
			cmdExit(0)
		},
	}

	persistentFlags := cmd.PersistentFlags()
	persistentFlags.SetInterspersed(false)

	persistentFlags.BoolVarP(&ignoreNewline, "ignore-newline", "l", true, "ignore leading a newline while sniffing the input")

	flags := cmd.Flags()
	flags.SetInterspersed(false)
	flags.BoolVarP(&passThrough, "pass-through", "p", false, "pass-through stdin to stdout")
	return cmd
}

var rootCmd = newRootCmd(builders)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		// fmt.Println(err)
		cmdExit(1)
	}
}

func init() {
}
