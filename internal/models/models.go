package models

type Chore struct {
	ID          int64
	Description string
	Points      int
	IsRequired  bool
	IsCompleted bool
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

type AssignmentDisplay struct {
	ID          int64
	ChildName   string
	ChoreName   string
	IsCompleted bool
}
