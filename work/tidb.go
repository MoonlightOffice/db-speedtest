package main

import (
	"crypto/tls"
	"database/sql"
	"fmt"
	"time"

	"github.com/go-sql-driver/mysql"
)

type TiDB struct {
	pool *sql.DB
}

func NewTiDB() (*TiDB, error) {
	mysql.RegisterTLSConfig("tidb", &tls.Config{
		MinVersion: tls.VersionTLS12,
		ServerName: "gateway01.ap-southeast-1.prod.aws.tidbcloud.com",
	})

	var err error
	pool, err := sql.Open("mysql", ConfigTiDBDSN)
	if err != nil {
		return nil, err
	}

	// Set up table
	stmt := `CREATE TABLE IF NOT EXISTS sample (
		id VARCHAR(50) PRIMARY KEY,
		name TEXT NOT NULL
	)`

	_, err = pool.Exec(stmt)
	if err != nil {
		return nil, err
	}

	return &TiDB{pool: pool}, nil
}

func (tidb TiDB) Close() error {
	return tidb.pool.Close()
}

func (tidb TiDB) DeleteAll() {
	stmt := `DELETE FROM sample`
	_, err := tidb.pool.Exec(stmt)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Delete: OK")
}

func (tidb TiDB) Run10Writes() {
	stmt := `INSERT INTO sample (id, name) VALUES (?, 'John Smith')`

	start := time.Now() // Count start

	for i := 0; i < 10; i++ {
		uniqueId := fmt.Sprintf("%d-%d", time.Now().UnixMilli(), i)
		_, err := tidb.pool.Exec(stmt, uniqueId)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	end := time.Now() // Count end

	fmt.Println("Write Result:", end.Sub(start))
}

func (tidb TiDB) Run10Reads() {
	stmt := `SELECT id, name FROM sample LIMIT 10`

	start := time.Now() // Count start

	rows, err := tidb.pool.Query(stmt)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer rows.Close()

	texts := make([]string, 0)
	for rows.Next() {
		var id, name string
		err = rows.Scan(&id, &name)
		if err != nil {
			fmt.Println(err)
			return
		}

		texts = append(texts, fmt.Sprintf("ID: %s, NAME: %s", id, name))
	}

	end := time.Now() // Count End

	for _, text := range texts {
		fmt.Println(text)
	}
	fmt.Println()
	fmt.Println("Read Result:", end.Sub(start))
}
