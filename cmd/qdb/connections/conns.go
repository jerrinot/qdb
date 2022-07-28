package connections

import (
	"github.com/spf13/cobra"
	"qdb/pkg/qdb"
)

var connectionName string
var connectionUrl string

var ConnectionsCmd = &cobra.Command{
	Use:     "conn",
	Aliases: []string{"shell"},
	Short:   "Manage Connections to QuestDB",
	Args:    cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		return qdb.ManageConnections()
	},
}

var listConnectionsCmd = &cobra.Command{
	Use:   "list",
	Short: "List Connections",
	Args:  cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		return qdb.ListConnections()
	},
}

var addConnectionCmd = &cobra.Command{
	Use:   "add",
	Short: "Add Connection",
	Args:  cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		return qdb.AddConnection(connectionName, connectionUrl)
	},
	Example: "qdb connections add --name localhost --url http://localhost:9000",
}

var deleteConnectionCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete Connection",
	Args:  cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		return qdb.DeleteConnection(connectionName)
	},
}

func init() {
	ConnectionsCmd.AddCommand(listConnectionsCmd)
	ConnectionsCmd.AddCommand(addConnectionCmd)
	ConnectionsCmd.AddCommand(deleteConnectionCmd)
	addConnectionCmd.Flags().StringVarP(&connectionName, "name", "n", "", "Connection name")
	addConnectionCmd.Flags().StringVarP(&connectionUrl, "url", "u", "", "Connection URL")
	_ = addConnectionCmd.MarkFlagRequired("name")
	_ = addConnectionCmd.MarkFlagRequired("url")

	deleteConnectionCmd.Flags().StringVarP(&connectionName, "name", "n", "", "Connection name")
	_ = deleteConnectionCmd.MarkFlagRequired("name")
}
