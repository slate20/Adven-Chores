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

// function to get all assignments from the database
func GetAllAssignments(db *sql.DB) ([]*Assignment, error) {
	var assignments []*Assignment

	rows, err := db.Query("SELECT id, child_id, chore_id, is_completed FROM assignments")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		assignment := &Assignment{}
		err := rows.Scan(&assignment.ID, &assignment.ChildID, &assignment.ChoreID, &assignment.IsCompleted)
		if err != nil {
			return nil, err
		}
		assignments = append(assignments, assignment)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return assignments, nil
}

// function to delete an assignment from the database
func DeleteAssignment(db *sql.DB, id int64) error {
	result, err := db.Exec("DELETE FROM assignments WHERE id = ?", id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no assignment found with ID %d", id)
	}

	return nil
}

// function to assign a chore to a child
func AssignChoreToChild(db *sql.DB, childID int64, choreID int64) error {
	// Check if the child and chore exist
	child, err := GetChildByID(db, childID)
	if err != nil {
		return fmt.Errorf("child not found: %v", err)
	}
	chore, err := GetChoreByID(db, choreID)
	if err != nil {
		return fmt.Errorf("chore not found: %v", err)
	}

	// Check if the assignment already exists
	var exists bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM assignments WHERE child_id = ? AND chore_id = ?)", child.ID, chore.ID).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check if assignment exists: %v", err)
	}
	if exists {
		return fmt.Errorf("chore %d is already assigned to child %d", chore.ID, child.ID)
	}

	// Create the assignment record
	assignment := &Assignment{
		ChildID:     child.ID,
		ChoreID:     chore.ID,
		IsCompleted: false,
	}

	err = assignment.Save(db)
	if err != nil {
		return fmt.Errorf("failed to assign chore to child: %v", err)
	}

	return nil
}

// function to unassign a chore from a child
func UnassignChoreFromChild(db *sql.DB, childID int64, choreID int64) error {
	// Check if the child and chore exist
	child, err := GetChildByID(db, childID)
	if err != nil {
		return fmt.Errorf("child not found: %v", err)
	}
	chore, err := GetChoreByID(db, choreID)
	if err != nil {
		return fmt.Errorf("chore not found: %v", err)
	}

	// Check if the assignment exists
	var exists bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM assignments WHERE child_id = ? AND chore_id = ?)", child.ID, chore.ID).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check if assignment exists: %v", err)
	}
	if !exists {
		return fmt.Errorf("chore %d is not assigned to child %d", chore.ID, child.ID)
	}

	// Delete the assignment record
	result, err := db.Exec("DELETE FROM assignments WHERE child_id = ? AND chore_id = ?", child.ID, chore.ID)
	if err != nil {
		return fmt.Errorf("failed to unassign chore from child: %v", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %v", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no assignment found for child %d and chore %d", child.ID, chore.ID)
	}

	return nil
}

// function to mark an assignment as completed
func CompleteAssignment(db *sql.DB, id int64) error {
	// Check if the assignment exists
	assignment, err := GetAssignmentByID(db, id)
	if err != nil {
		return fmt.Errorf("assignment not found: %v", err)
	}

	// Check if the assignment is already completed
	if assignment.IsCompleted {
		return fmt.Errorf("assignment is already completed")
	}

	// Get the chore details
	chore, err := GetChoreByID(db, assignment.ChoreID)
	if err != nil {
		return fmt.Errorf("failed to get chore details: %v", err)
	}

	// Update the assignment record
	assignment.IsCompleted = true
	err = assignment.Save(db)
	if err != nil {
		return fmt.Errorf("failed to mark assignment as completed: %v", err)
	}

	// Update the child's points
	child, err := GetChildByID(db, assignment.ChildID)
	if err != nil {
		return fmt.Errorf("failed to get child details: %v", err)
	}
	child.Points += chore.Points
	err = child.Save(db)
	if err != nil {
		return fmt.Errorf("failed to update child points: %v", err)
	}

	return nil
}

// function to unmark an assignment as completed
func UncompleteAssignment(db *sql.DB, id int64) error {
	// Check if the assignment exists
	assignment, err := GetAssignmentByID(db, id)
	if err != nil {
		return fmt.Errorf("assignment not found: %v", err)
	}

	// Update the assignment record
	assignment.IsCompleted = false
	err = assignment.Save(db)
	if err != nil {
		return fmt.Errorf("failed to unmark assignment as completed: %v", err)
	}

	return nil
}

// function to get all assignments for a child
func GetChoresByChild(db *sql.DB, childID int64) ([]*Chore, error) {
	// Check if the child exists
	_, err := GetChildByID(db, childID)
	if err != nil {
		return nil, fmt.Errorf("child not found: %v", err)
	}

	// join the assignments and chores tables
	query := `
		SELECT c.id, c.description, c.points, c.is_required, c.due_date, a.is_completed
		FROM chores c
		JOIN assignments a ON c.id = a.chore_id
		WHERE a.child_id = ?
	`

	rows, err := db.Query(query, childID)
	if err != nil {
		return nil, fmt.Errorf("failed to get chores by child: %v", err)
	}
	defer rows.Close()

	var chores []*Chore
	for rows.Next() {
		chore := &Chore{}
		var isCompleted bool
		err := rows.Scan(&chore.ID, &chore.Description, &chore.Points, &chore.IsRequired, &chore.DueDate, &isCompleted)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}
		chore.IsCompleted = isCompleted
		chores = append(chores, chore)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %v", err)
	}

	return chores, nil
}
