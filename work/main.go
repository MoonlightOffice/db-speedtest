package main

import (
	"fmt"
)

const (
	// CockroachDB
	ConfigCRDBConn string = "postgresql://<YOUR_CONNECTION_STRING>"

	// TiDB
	ConfigTiDBServerName string = "gateway01.ap-southeast-1.prod.aws.tidbcloud.com"
	ConfigTiDBDSN        string = "YOUR_DSN"
)

func main() {
	//db, err := NewCRDB()
	db, err := NewTiDB()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	db.Run10Writes()
	//db.Run10Reads()
	//db.DeleteAll()
}
