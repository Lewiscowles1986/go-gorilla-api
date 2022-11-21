package settings

import (
	"fmt"
)

func GetDBConnStr(connType, username, password, database string) string {
	if connType == "postgres" {
		return fmt.Sprintf(
			"user=%s password=%s dbname=%s sslmode=%s",
			username, password, database, Getenv("DB_SSL_MODE", "disable"))
	}
	if connType == "sqlite3" && database != ":memory:" {
		return fmt.Sprintf("%s.db", database)
	}
	return database
}
