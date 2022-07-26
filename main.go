package main

import (
	"qdb/cmd/qdb"
)

func main() {
	//err := qdb.LoadConfig(qdb.SqlCmd)
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println(qdb.DefaultProfile)
	//err = qdb.SaveDefaultProfile(qdb.DefaultProfile + ".")
	//if err != nil {
	//	panic(err)
	//}
	qdb.Execute()
}
