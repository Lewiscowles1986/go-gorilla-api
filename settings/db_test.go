package settings

import (
	"testing"
)

func TestSqlite3InMemory(t *testing.T) {
	result := GetDBConnStr("sqlite3", "", "", ":memory:")
	if result != ":memory:" {
		t.Errorf("SQLite in-memory resolution failed")
	}
}

func TestSqlite3AnyOtherFile(t *testing.T) {
	expected := "someOtherDbName.db"
	result := GetDBConnStr("sqlite3", "", "", "someOtherDbName")
	if result != expected {
		t.Errorf("SQLite database file result incorrect."+
			"Expected: ''%s'. Received: '%s'", expected, result)
	}
}

func TestPostGresDB(t *testing.T) {
	expected := "user=username password=password dbname=someOtherDbName sslmode=disable"
	result := GetDBConnStr("postgres", "username", "password", "someOtherDbName")
	if result != expected {
		t.Errorf("postgres connection string resolution failed. "+
			"Expected: '%s'. Received: '%s'", expected, result)
	}
}
