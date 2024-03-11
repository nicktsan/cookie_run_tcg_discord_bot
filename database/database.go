package db

import (
	"database/sql"
	"fmt"
)

func SqlConfig(ConnectionStr string) (*sql.DB, error) {
	conn, err := sql.Open("postgres", ConnectionStr)
	if err != nil {
		return nil, err
	}
	fmt.Println("DB Connection successful.")
	return conn, nil
}
