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

	// Create an Adven-Chores directory in the user's home directory
	dbDir := filepath.Join(homeDir, "Adven-Chores")
	err = os.MkdirAll(dbDir, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	// Create the database path
	dbPath := filepath.Join(dbDir, "advenchores.db")
	log.Printf("Database path: %s", dbPath)

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

	createUsersTable := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT UNIQUE NOT NULL,
		email TEXT UNIQUE NOT NULL,
		password_hash TEXT NOT NULL,
		parent_pin INTEGER NOT NULL DEFAULT 1234,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	_, err := DB.Exec(createUsersTable)
	if err != nil {
		log.Fatal(err)
	}

	createPasswordResetTokensTable := `
	CREATE TABLE IF NOT EXISTS password_reset_tokens (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		token TEXT NOT NULL,
		FOREIGN KEY (user_id) REFERENCES users(id)
	);`

	_, err = DB.Exec(createPasswordResetTokensTable)
	if err != nil {
		log.Fatal(err)
	}

	createSecurityQuestionsTable := `
	CREATE TABLE IF NOT EXISTS security_questions (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		question TEXT NOT NULL,
		answer_hash TEXT NOT NULL,
		FOREIGN KEY (user_id) REFERENCES users(id)
	);`

	_, err = DB.Exec(createSecurityQuestionsTable)
	if err != nil {
		log.Fatal(err)
	}

	createChoresTable := `
	CREATE TABLE IF NOT EXISTS chores (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		description TEXT NOT NULL,
		points INTEGER NOT NULL,
		is_required BOOLEAN NOT NULL,
		is_completed BOOLEAN NOT NULL DEFAULT 0
	);`

	_, err = DB.Exec(createChoresTable)
	if err != nil {
		log.Fatal(err)
	}

	createChildrenTable := `
	CREATE TABLE IF NOT EXISTS children (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		name TEXT NOT NULL,
		job TEXT NOT NULL,
		rewards STRING,
		points INTEGER DEFAULT 0
	);`

	_, err = DB.Exec(createChildrenTable)
	if err != nil {
		log.Fatal(err)
	}

	createAssignmentsTable := `
	CREATE TABLE IF NOT EXISTS assignments (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
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
		user_id INTEGER NOT NULL,
		description TEXT NOT NULL,
		point_cost INTEGER NOT NULL
	);`

	_, err = DB.Exec(createRewardsTable)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Database tables initialized")
}
