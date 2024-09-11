package models

import (
	"database/sql"
	"fmt"
)

// TODO: Add username and email update functions
func UpdateUsername(db *sql.DB, id int, username string) error {
	_, err := db.Exec("UPDATE users SET username = ? WHERE id = ?", username, id)
	if err != nil {
		return err
	}
	return nil
}

// function to get user details from the database by ID
func GetUserByID(db *sql.DB, id int) (*User, error) {
	query := "SELECT id, username, email, parent_pin FROM users WHERE id = ?"
	row := db.QueryRow(query, id)

	user := &User{}
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.ParentPin)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no user found with ID %d", id)
		}
		return nil, err
	}

	return user, nil
}
