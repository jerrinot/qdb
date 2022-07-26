package qdb

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var version = "0.0.1"

var rootCmd = &cobra.Command{
	Use:     "qdb",
	Short:   "qdb - a simple CLI to mess with QuestDB",
	Version: version,
	Long: `qdb is a super fancy CLI for QuestDB
One can use qdb to modify or inspect QuestDB straight from the terminal`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Whoops. There was an error while executing your CLI '%s'", err)
		os.Exit(1)
	}
}
