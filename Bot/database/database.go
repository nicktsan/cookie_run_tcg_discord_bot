package database

import (
	"database/sql"
	"log"

	errFunc "discordbot/cookieruntcg_bot/error"
)

func ConnectDB(dbType string, ConnectionStr string) *sql.DB {
	conn, err := sql.Open(dbType, ConnectionStr)
	errFunc.CheckNilErrPanic("Error occured while attempting to connect to DB.", err)
	log.Println("DB Connection successful.")
	return conn
}
