package models

type User struct {
	ID        int
	Username  string
	Email     string
	ParentPin int
}

type Chore struct {
	ID          int64
	UserID      int
	Description string
	Points      int
	IsRequired  bool
	IsCompleted bool
}

type Child struct {
	ID      int64
	UserID  int
	Name    string
	Job     string
	Points  int
	Rewards string
}
type Reward struct {
	ID          int64
	UserID      int
	Description string
	PointCost   int
}

type Assignment struct {
	ID          int64
	ChoreID     int64
	ChildID     int64
	ChildName   string
	IsCompleted bool
	Chore       *Chore
}
