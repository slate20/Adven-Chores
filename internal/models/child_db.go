package models

import (
	"database/sql"
	"fmt"
	"log"
)

// function to save a child to the database
func (c *Child) Save(db *sql.DB) error {
	// If the child is new, insert it
	if c.ID == 0 {
		result, err := db.Exec("INSERT INTO children (user_id, name, job, points, rewards) VALUES (?, ?, ?, ?, ?)",
			c.UserID, c.Name, c.Job, c.Points, c.Rewards)
		if err != nil {
			return err
		}

		c.ID, err = result.LastInsertId()
		if err != nil {
			return err
		}
	} else {
		// If the child is not new, update it
		log.Printf("Updating child %d for user %d", c.ID, c.UserID)
		result, err := db.Exec("UPDATE children SET name = ?, job = ?, points = ?, rewards = ? WHERE id = ? AND user_id = ?",
			c.Name, c.Job, c.Points, c.Rewards, c.ID, c.UserID)
		if err != nil {
			log.Printf("Failed to update child %d: %v", c.ID, err)
			return err
		}
		rowsAffected, err := result.RowsAffected()
		if err != nil {
			log.Printf("Failed to get number of rows affected for child %d: %v", c.ID, err)
			return err
		}
		if rowsAffected == 0 {
			log.Printf("No rows affected for child %d", c.ID)
			return fmt.Errorf("no row was updated for child %d", c.ID)
		}
	}
	return nil
}

// function to get a child by ID from the database
func GetChildByID(db *sql.DB, id int64) (*Child, error) {
	query := "SELECT id, user_id, name, job, points, rewards FROM children WHERE id = ?"
	row := db.QueryRow(query, id)

	child := &Child{}
	err := row.Scan(&child.ID, &child.UserID, &child.Name, &child.Job, &child.Points, &child.Rewards)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no child found with ID %d", id)
		}
		return nil, err
	}

	return child, nil
}

// function to get all children from the database
func GetAllChildren(db *sql.DB) ([]*Child, error) {
	var children []*Child

	rows, err := db.Query("SELECT id, name, job, points, rewards FROM children")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		child := &Child{}
		err := rows.Scan(&child.ID, &child.Name, &child.Job, &child.Points, &child.Rewards)
		if err != nil {
			return nil, err
		}
		children = append(children, child)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return children, nil
}

// function to get children by user ID from the database
func GetChildrenByUserID(db *sql.DB, userID int) ([]*Child, error) {
	var children []*Child

	rows, err := db.Query("SELECT id, name, job, points, rewards FROM children WHERE user_id = ?", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		child := &Child{}
		err := rows.Scan(&child.ID, &child.Name, &child.Job, &child.Points, &child.Rewards)
		if err != nil {
			return nil, err
		}
		children = append(children, child)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return children, nil
}

// function to delete a child from the database
func DeleteChild(db *sql.DB, id int64) error {
	result, err := db.Exec("DELETE FROM children WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("error deleting child: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no child found with ID %d", id)
	}

	return nil
}
