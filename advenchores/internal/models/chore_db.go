package models

import (
	"database/sql"
	"fmt"
	"log"
)

// function to save a chore to the database
func (c *Chore) Save(db *sql.DB) error {
	// If the chore is new, insert it
	if c.ID == 0 {
		result, err := db.Exec("INSERT INTO chores (description, points, is_required) VALUES (?, ?, ?)",
			c.Description, c.Points, c.IsRequired)
		if err != nil {
			return err
		}

		c.ID, err = result.LastInsertId()
		if err != nil {
			return err
		}
	} else {
		// If the chore is not new, update it
		_, err := db.Exec("UPDATE chores SET description = ?, points = ?, is_required = ? WHERE id = ?",
			c.Description, c.Points, c.IsRequired, c.ID)
		if err != nil {
			return err
		}
	}
	return nil
}

// function to get a chore by ID from the database
func GetChoreByID(db *sql.DB, id int64) (*Chore, error) {
	query := "SELECT id, description, points, is_required FROM chores WHERE id = ?"
	row := db.QueryRow(query, id)

	chore := &Chore{}
	err := row.Scan(&chore.ID, &chore.Description, &chore.Points, &chore.IsRequired)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no chore found with ID %d", id)
		}
		return nil, err
	}

	return chore, nil
}

// function to get all chores from the database
func GetAllChores(db *sql.DB) ([]*Chore, error) {
	var chores []*Chore

	log.Printf("Getting all chores from database")

	rows, err := db.Query("SELECT id, description, points, is_required FROM chores")
	if err != nil {
		log.Printf("Error getting all chores from database: %v", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		chore := &Chore{}
		err := rows.Scan(&chore.ID, &chore.Description, &chore.Points, &chore.IsRequired)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			return nil, err
		}

		log.Printf("Found chore: %v", chore)

		chores = append(chores, chore)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error iterating over rows: %v", err)
		return nil, err
	}

	log.Printf("Found %d chores", len(chores))

	return chores, nil
}

// function to delete a chore from the database
func DeleteChore(db *sql.DB, id int64) error {
	result, err := db.Exec("DELETE FROM chores WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("error deleting chore: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no chore found with ID %d", id)
	}

	return nil
}
