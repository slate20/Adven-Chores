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
			a.ChildID, a.Chore.ID, a.IsCompleted)
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
			a.ChildID, a.Chore.ID, a.IsCompleted, a.ID)
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

	assignment := &Assignment{Chore: &Chore{}}
	err := row.Scan(&assignment.ID, &assignment.ChildID, &assignment.Chore.ID, &assignment.IsCompleted)
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
	query := `
		SELECT a.id, a.child_id, ch.name, c.id, c.description, c.points, c.is_required, a.is_completed
		FROM assignments a
		JOIN chores c ON a.chore_id = c.id
		JOIN children ch ON a.child_id = ch.id
		`

	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get assignments: %v", err)
	}
	defer rows.Close()

	var assignments []*Assignment
	for rows.Next() {
		a := &Assignment{Chore: &Chore{}}
		err := rows.Scan(
			&a.ID,
			&a.ChildID,
			&a.ChildName,
			&a.Chore.ID,
			&a.Chore.Description,
			&a.Chore.Points,
			&a.Chore.IsRequired,
			&a.IsCompleted,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}
		assignments = append(assignments, a)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate rows: %v", err)
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
		Chore:       chore,
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

	// Update the assignment record
	assignment.IsCompleted = true
	err = assignment.Save(db)
	if err != nil {
		return fmt.Errorf("failed to mark assignment as completed: %v", err)
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

// function to reward and remove an assignment once it has been completed
func RewardAssignment(db *sql.DB, id int64) error {
	// Check if the assignment exists
	assignment, err := GetAssignmentByID(db, id)
	if err != nil {
		return fmt.Errorf("assignment not found: %v", err)
	}

	// Check if the assignment is completed
	if !assignment.IsCompleted {
		return fmt.Errorf("assignment is not completed")
	}

	// Get the chore details
	chore, err := GetChoreByID(db, assignment.Chore.ID)
	if err != nil {
		return fmt.Errorf("failed to get chore details: %v", err)
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

	// call DeleteAssignment() to remove the assignment from the database
	err = DeleteAssignment(db, id)
	if err != nil {
		return fmt.Errorf("failed to delete assignment: %v", err)
	}

	return nil
}

// function to get all assignments for a child
func GetAssignmentsByChild(db *sql.DB, childID int64) ([]*Assignment, error) {
	// Check if the child exists
	_, err := GetChildByID(db, childID)
	if err != nil {
		return nil, fmt.Errorf("child not found: %v", err)
	}

	// join the assignments and chores tables
	query := `
		SELECT a.id, a.child_id, ch.name, a.chore_id, a.is_completed, c.description, c.points, c.is_required
		FROM assignments a
		JOIN chores c ON a.chore_id = c.id
		JOIN children ch ON a.child_id = ch.id
		WHERE a.child_id = ?
	`

	rows, err := db.Query(query, childID)
	if err != nil {
		return nil, fmt.Errorf("failed to get assignments by child: %v", err)
	}
	defer rows.Close()

	var assignments []*Assignment
	for rows.Next() {
		assignment := &Assignment{Chore: &Chore{}}
		err := rows.Scan(
			&assignment.ID,
			&assignment.ChildID,
			&assignment.ChildName,
			&assignment.Chore.ID,
			&assignment.IsCompleted,
			&assignment.Chore.Description,
			&assignment.Chore.Points,
			&assignment.Chore.IsRequired,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}
		assignments = append(assignments, assignment)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %v", err)
	}

	return assignments, nil
}
