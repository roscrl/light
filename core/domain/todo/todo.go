package todo

import "time"

type Status string

const (
	Pending Status = "pending"
	Done    Status = "done"
)

type Todo struct {
	ID        string
	Task      string
	Status    Status
	CreatedAt time.Time
}
