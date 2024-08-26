package models

import (
	"time"
)

type Chore struct {
	ID          int64
	Description string
	Points      int
	IsRequired  bool
	DueDate     time.Time
}

type Child struct {
	ID     int64
	Name   string
	Job    string
	Points int
}

type Reward struct {
	ID          int64
	Description string
	PointCost   int
}

type Assignment struct {
	ID          int64
	ChildID     int64
	ChoreID     int64
	IsCompleted bool
}
