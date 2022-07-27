package qdb

var Query string
var Profile string

func init() {
	SqlCmd.Flags().StringVarP(&Query, "query", "q", "", "Query to run in non-interactive mode")
	SqlCmd.Flags().StringVarP(&Profile, "profile", "p", "", "QuestDB profile selection")
	rootCmd.AddCommand(SqlCmd)
	rootCmd.AddCommand(ConnectionsCmd)
}
