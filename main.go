package main

import (
	"flag"
	"os"
	"time"

	"./settings"
)

func main() {
	var wait time.Duration
	flag.DurationVar(&wait, "graceful-timeout", time.Minute*1, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.Parse()

	a := App{}
	a.Initialize(
		settings.Getenv("APP_DB_TYPE", "sqlite3"),
		settings.GetDBConnStr(
			settings.Getenv("APP_DB_TYPE", "sqlite3"),
			os.Getenv("APP_DB_USERNAME"),
			os.Getenv("APP_DB_PASSWORD"),
			settings.Getenv("APP_DB_NAME", "database")))

	a.Run(":8080", wait)
}
