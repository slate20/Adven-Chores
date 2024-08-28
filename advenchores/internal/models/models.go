package models

type Chore struct {
	ID          int64
	Description string
	Points      int
	IsRequired  bool
	IsCompleted bool
}

type Child struct {
	ID      int64
	Name    string
	Job     string
	Points  int
	Rewards string
}
type Reward struct {
	ID          int64
	Description string
	PointCost   int
}

type Assignment struct {
	ID          int64
	ChildID     int64
	ChildName   string
	IsCompleted bool
	Chore       *Chore
}
