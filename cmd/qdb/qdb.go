package qdb

import (
	"qdb/pkg/qdb"

	"github.com/spf13/cobra"
)

var sqlCmd = &cobra.Command{
	Use:     "sql",
	Aliases: []string{"shell"},
	Short:   "Run SQL shell",
	Args:    cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		qdb.RunSqlShell()
	},
}

func init() {
	rootCmd.AddCommand(sqlCmd)
}
