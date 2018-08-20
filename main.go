package main

import (
	"os"

	"./settings"
)

func main() {
	a := App{}
	a.Initialize(
		settings.Getenv("APP_DB_TYPE", "sqlite3"),
		settings.GetDBConnStr(
			settings.Getenv("APP_DB_TYPE", "sqlite3"),
			os.Getenv("APP_DB_USERNAME"),
			os.Getenv("APP_DB_PASSWORD"),
			settings.Getenv("APP_DB_NAME", "database")))

	a.Run(":8080")
}
