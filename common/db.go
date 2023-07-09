package common

import (
	"database/sql"
	"testing"

	_ "embed"

	_ "github.com/mattn/go-sqlite3"
)

//go:embed schema.sql
var schema string

//go:embed grail.sql
var standardData string

// Open connects to a sqlite database and loads the schema. If the database
// file does not exist, it will be created.
func Open(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	// Load the schema
	_, err = db.Exec(schema)
	if err != nil {
		return nil, err
	}

	return db, nil
}

// Populate loads some starter data into the database.
func Populate(db *sql.DB) error {
	_, err := db.Exec(standardData)
	return err
}

// TestDB creates a new in-memory database for testing. The schema is loaded
// and some test data is populated. If there is an error, the test is aborted
// (t.Fatal).
func TestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %s", err)
	}

	// Load the schema
	_, err = db.Exec(schema)
	if err != nil {
		db.Close()
		t.Fatalf("Error loading schema: %s", err)
	}

	// Load the test data
	err = Populate(db)
	if err != nil {
		db.Close()
		t.Fatalf("Error loading test data: %s", err)
	}

	return db
}
