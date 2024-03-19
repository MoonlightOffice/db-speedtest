package main

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type CRDB struct {
	pool *sql.DB
}

func NewCRDB() (*CRDB, error) {
	var err error
	pool, err := sql.Open("pgx", ConfigCRDBConn)
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

	return &CRDB{pool: pool}, nil
}

func (crdb CRDB) Close() error {
	return crdb.pool.Close()
}

func (crdb CRDB) DeleteAll() {
	stmt := `DELETE FROM sample`
	_, err := crdb.pool.Exec(stmt)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Delete: OK")
}

func (crdb CRDB) Run10Writes() {
	stmt := `INSERT INTO sample (id, name) VALUES ($1, 'John Smith')`

	start := time.Now() // Count start

	for i := 0; i < 10; i++ {
		uniqueId := fmt.Sprintf("%d-%d", time.Now().UnixMilli(), i)
		_, err := crdb.pool.Exec(stmt, uniqueId)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	end := time.Now() // Count end

	fmt.Println("Write Result:", end.Sub(start))
}

func (crdb CRDB) Run10Reads() {
	stmt := `SELECT id, name FROM sample LIMIT 10`

	start := time.Now() // Count start

	rows, err := crdb.pool.Query(stmt)
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
