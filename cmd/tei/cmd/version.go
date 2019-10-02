package cmd

import (
	"fmt"
	"log"
	"runtime"
	"text/template"

	"github.com/spf13/cobra"
)

func newVersionCmd() *cobra.Command {
	// VersionCmd represents the file command
	cmd := &cobra.Command{
		Use:                   "version",
		DisableFlagsInUseLine: true,
		Short:                 "Show " + cmdName + " version inforation",
		Long:                  "",
		Args:                  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			if cmdVer != "" {
				t := template.Must(template.New("footer").Parse(`{{define "footer"}}
{{- if .Name}}{{printf "%s" .Name}}{{end}}
{{printf " version:    %s" .Version}}
{{printf " go version: %s" .Gover}}
{{printf " os/arch:    %s/%s" .Goos .Goarch}}
{{- if .Hash}}{{printf "\n git commit: %s" .Hash}}{{end}}
{{end}}`))
				err := t.ExecuteTemplate(cmd.OutOrStdout(), "footer", struct {
					Name    string
					Version string
					Gover   string
					Goos    string
					Goarch  string
					Hash    string
				}{
					Name:    cmdName,
					Gover:   runtime.Version(),
					Version: cmdVer,
					Goos:    runtime.GOOS,
					Goarch:  runtime.GOARCH,
					Hash:    cmdCommitHash,
				})
				if err != nil {
					log.Fatal(err)
				}
			} else {
				fmt.Fprintln(cmd.OutOrStdout(), "built with no version infromation.")
			}
		},
	}

	flags := cmd.Flags()
	flags.SetInterspersed(false)

	return cmd
}

func init() {
	rootCmd.AddCommand(newVersionCmd())
}
