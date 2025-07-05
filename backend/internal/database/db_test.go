package database

import (
	"testing"

	_ "modernc.org/sqlite"
)

func TestInitDB(t *testing.T) {
	// Use in-memory database for testing
	InitDB("file::memory:?cache=shared")

	tables := []string{"users", "products", "orders"}

	for _, table := range tables {
		var name string
		// SQLite stores table names in lower case
		err := DB.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name=?", table).Scan(&name)
		if err != nil {
			t.Fatalf("Table '%s' was not created: %v", table, err)
		}
		if name != table {
			t.Fatalf("Expected table name '%s', but got '%s'", table, name)
		}
	}
}