package sql

import (
	"github.com/spf13/cobra"
	"qdb/pkg/qdb"
)

var Query string
var ConnectionName string

var SqlCmd = &cobra.Command{
	Use:     "sql",
	Aliases: []string{"shell"},
	Short:   "Run SQL shell",
	Args:    cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		return qdb.RunSqlShell(Query, ConnectionName)
	},
}

func init() {
	SqlCmd.Flags().StringVarP(&Query, "query", "q", "", "Query to run in non-interactive mode")
	SqlCmd.Flags().StringVarP(&ConnectionName, "connection", "c", "", "QuestDB connection selection")
}
