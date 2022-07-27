package qdb

import (
	"github.com/spf13/cobra"
	"qdb/pkg/qdb"
)

var version = "0.0.1"

var rootCmd = &cobra.Command{
	Use:     "qdb",
	Short:   "qdb - a simple CLI to mess with QuestDB",
	Version: version,
	Long: `qdb is a super fancy CLI for QuestDB
One can use qdb to modify or inspect QuestDB straight from the terminal`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return qdb.LoadConfig(cmd)
	},
	PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
		return qdb.SaveConfig()
	},
}

func Execute() error {
	return rootCmd.Execute()
}
