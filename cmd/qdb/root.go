package qdb

import (
	"github.com/spf13/cobra"
	"qdb/cmd/qdb/connections"
	"qdb/cmd/qdb/sql"
	"qdb/pkg/qdb"
)

var version = "0.0.1"

var rootCmd = &cobra.Command{
	Use:     "qdb",
	Short:   "qdb - a simple CLI to mess with QuestDB",
	Version: version,
	Long: `qdb is a super fancy CLI for QuestDB
One can use qdb to modify or inspect QuestDB straight from the terminal`,
	PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
		return qdb.SaveConfig()
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	_ = qdb.LoadConfig()
	rootCmd.AddCommand(sql.SqlCmd)
	rootCmd.AddCommand(connections.ConnectionsCmd)
}
