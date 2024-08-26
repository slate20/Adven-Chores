package models

import (
	"database/sql"
	"fmt"
)

// function to save a child to the database
func (c *Child) Save(db *sql.DB) error {
	// If the child is new, insert it
	if c.ID == 0 {
		result, err := db.Exec("INSERT INTO children (name, job, points) VALUES (?, ?, ?)",
			c.Name, c.Job, c.Points)
		if err != nil {
			return err
		}

		c.ID, err = result.LastInsertId()
		if err != nil {
			return err
		}
	} else {
		// If the child is not new, update it
		_, err := db.Exec("UPDATE children SET name = ?, job = ?, points = ? WHERE id = ?",
			c.Name, c.Job, c.Points, c.ID)
		if err != nil {
			return err
		}
	}
	return nil
}

// function to get a child by ID from the database
func GetChildByID(db *sql.DB, id int64) (*Child, error) {
	query := "SELECT id, name, job, points FROM children WHERE id = ?"
	row := db.QueryRow(query, id)

	child := &Child{}
	err := row.Scan(&child.ID, &child.Name, &child.Job, &child.Points)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no child found with ID %d", id)
		}
		return nil, err
	}

	return child, nil
}
