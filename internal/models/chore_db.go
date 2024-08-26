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
func GetAllChores(db *sql.DB) ([]*Chore, error) {
	var chores []*Chore

	rows, err := db.Query("SELECT id, description, points, is_required, due_date FROM chores")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		chore := &Chore{}
		err := rows.Scan(&chore.ID, &chore.Description, &chore.Points, &chore.IsRequired, &chore.DueDate)
		if err != nil {
			return nil, err
		}

		chores = append(chores, chore)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

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
