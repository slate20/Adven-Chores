package models

import (
	"database/sql"
	"fmt"
)

// function to save an assignment to the database
func (a *Assignment) Save(db *sql.DB) error {
	// If the assignment is new, insert it
	if a.ID == 0 {
		result, err := db.Exec("INSERT INTO assignments (child_id, chore_id, is_completed) VALUES (?, ?, ?)",
			a.ChildID, a.ChoreID, a.IsCompleted)
		if err != nil {
			return err
		}

		a.ID, err = result.LastInsertId()
		if err != nil {
			return err
		}
	} else {
		// If the assignment is not new, update it
		_, err := db.Exec("UPDATE assignments SET child_id = ?, chore_id = ?, is_completed = ? WHERE id = ?",
			a.ChildID, a.ChoreID, a.IsCompleted, a.ID)
		if err != nil {
			return err
		}
	}
	return nil
}

// function to retrieve an assignment by ID from the database
func GetAssignmentByID(db *sql.DB, id int64) (*Assignment, error) {
	query := "SELECT id, child_id, chore_id, is_completed FROM assignments WHERE id = ?"
	row := db.QueryRow(query, id)

	assignment := &Assignment{}
	err := row.Scan(&assignment.ID, &assignment.ChildID, &assignment.ChoreID, &assignment.IsCompleted)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no assignment found with ID %d", id)
		}
		return nil, err
	}

	return assignment, nil
}
