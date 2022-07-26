package qdb

import (
	"qdb/pkg/qdb"

	"github.com/spf13/cobra"
)

var Query string

var sqlCmd = &cobra.Command{
	Use:     "sql",
	Aliases: []string{"shell"},
	Short:   "Run SQL shell",
	Args:    cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		return qdb.RunSqlShell(Query)
	},
}

func init() {
	sqlCmd.Flags().StringVarP(&Query, "query", "q", "", "Query to run in non-interactive mode")
	rootCmd.AddCommand(sqlCmd)
}
