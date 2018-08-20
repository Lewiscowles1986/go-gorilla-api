package settings

import (
	"fmt"
)

func GetDBConnStr(connType, username, password, database string) string {
	if connType == "postgresql" {
		return fmt.Sprintf(
			"user=%s password=%s dbname=%s",
			username, password, database)
	}
	if connType == "sqlite3" && database != ":memory:" {
		return fmt.Sprintf("%s.db", database)
	}
	return database
}
