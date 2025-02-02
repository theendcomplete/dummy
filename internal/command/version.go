package command

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/go-dummy/dummy/internal/exitcode"
)

func (e *Executor) initVersion() {
	e.versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Version",
		Run: func(_ *cobra.Command, _ []string) {
			fmt.Printf("Dummy version %s\n", e.version)
			os.Exit(exitcode.Success)
		},
	}

	e.rootCmd.AddCommand(e.versionCmd)
}
