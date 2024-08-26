package models

import (
	"database/sql"
	"fmt"
)

// function to save a chore to the database
func (c *Chore) Save(db *sql.DB) error {
	// If the chore is new, insert it
	if c.ID == 0 {
		result, err := db.Exec("INSERT INTO chores (description, points, is_required, due_date) VALUES (?, ?, ?, ?)",
			c.Description, c.Points, c.IsRequired, c.DueDate)
		if err != nil {
			return err
		}

		c.ID, err = result.LastInsertId()
		if err != nil {
			return err
		}
	} else {
		// If the chore is not new, update it
		_, err := db.Exec("UPDATE chores SET description = ?, points = ?, is_required = ?, due_date = ? WHERE id = ?",
			c.Description, c.Points, c.IsRequired, c.DueDate, c.ID)
		if err != nil {
			return err
		}
	}
	return nil
}

// function to get a chore by ID from the database
func GetChoreByID(db *sql.DB, id int64) (*Chore, error) {
	query := "SELECT id, description, points, is_required, due_date FROM chores WHERE id = ?"
	row := db.QueryRow(query, id)

	chore := &Chore{}
	err := row.Scan(&chore.ID, &chore.Description, &chore.Points, &chore.IsRequired, &chore.DueDate)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no chore found with ID %d", id)
		}
		return nil, err
	}

	return chore, nil
}

// function to get all chores from the database
func GetChores(db *sql.DB) ([]*Chore, error) {
	query := "SELECT id, description, points, is_required, due_date FROM chores"
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var chores []*Chore

	for rows.Next() {
		chore := &Chore{}
		err := rows.Scan(&chore.ID, &chore.Description, &chore.Points, &chore.IsRequired, &chore.DueDate)
		if err != nil {
			return nil, err
		}
		chores = append(chores, chore)
	}

	return chores, nil
}
