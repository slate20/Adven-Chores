package database

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB() (*sql.DB, error) {
	// Get the user's home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	// Create a ChoreQuest directory in the user's home directory
	dbDir := filepath.Join(homeDir, "ChoreQuest")
	err = os.MkdirAll(dbDir, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	// Create the database path
	dbPath := filepath.Join(dbDir, "chorequest.db")

	// Open the database
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal(err)
	}

	// Test the connection
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	DB = db
	log.Println("Database connection established")

	// Initialize the database
	initTables()

	return db, nil
}

func initTables() {
	// Create the tables
	createChoresTable := `
	CREATE TABLE IF NOT EXISTS chores (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		description TEXT NOT NULL,
		points INTEGER NOT NULL,
		is_required BOOLEAN NOT NULL,
		due_date TEXT
	);`

	_, err := DB.Exec(createChoresTable)
	if err != nil {
		log.Fatal(err)
	}

	createChildrenTable := `
	CREATE TABLE IF NOT EXISTS children (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		job TEXT NOT NULL,
		points INTEGER DEFAULT 0
	);`

	_, err = DB.Exec(createChildrenTable)
	if err != nil {
		log.Fatal(err)
	}

	createAssignmentsTable := `
	CREATE TABLE IF NOT EXISTS assignments (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		chore_id INTEGER,
		child_id INTEGER,
		is_completed BOOLEAN,
		FOREIGN KEY(chore_id) REFERENCES chores(id),
		FOREIGN KEY(child_id) REFERENCES children(id)
	);`

	_, err = DB.Exec(createAssignmentsTable)
	if err != nil {
		log.Fatal(err)
	}

	createRewardsTable := `
	CREATE TABLE IF NOT EXISTS rewards (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		description TEXT NOT NULL,
		point_cost INTEGER NOT NULL
	);`

	_, err = DB.Exec(createRewardsTable)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Database tables initialized")
}
