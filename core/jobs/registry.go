package jobs

import "github.com/roscrl/light/core/jobs/tododelete"

type JobName string

type JobFunc func(args map[string]any) error

type Registry map[JobName]JobFunc

const (
	TodoDelete JobName = "todo_delete"
)

func DefaultRegistry() Registry {
	return map[JobName]JobFunc{
		TodoDelete: tododelete.Run,
	}
}
