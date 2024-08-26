package models

import (
	"database/sql"
	"fmt"
)

// function to save a reward to the database
func (r *Reward) Save(db *sql.DB) error {
	// If the reward is new, insert it
	if r.ID == 0 {
		result, err := db.Exec("INSERT INTO rewards (description, point_cost) VALUES (?, ?)",
			r.Description, r.PointCost)
		if err != nil {
			return err
		}

		r.ID, err = result.LastInsertId()
		if err != nil {
			return err
		}
	} else {
		// If the reward is not new, update it
		_, err := db.Exec("UPDATE rewards SET description = ?, point_cost = ? WHERE id = ?",
			r.Description, r.PointCost, r.ID)
		if err != nil {
			return err
		}
	}
	return nil
}

// function to get a reward by ID from the database
func GetRewardByID(db *sql.DB, id int64) (*Reward, error) {
	query := "SELECT id, description, point_cost FROM rewards WHERE id = ?"
	row := db.QueryRow(query, id)

	reward := &Reward{}
	err := row.Scan(&reward.ID, &reward.Description, &reward.PointCost)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no reward found with ID %d", id)
		}
		return nil, err
	}

	return reward, nil
}
