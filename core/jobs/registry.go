package jobs

import "github.com/roscrl/light/core/jobs/tododelete"

type JobName string

type JobFunc func(args map[string]any) error

type JobNameToJobFuncRegistry map[JobName]JobFunc

const (
	TodoDelete JobName = "todo_delete"
)

func DefaultRegistry() JobNameToJobFuncRegistry {
	return map[JobName]JobFunc{
		TodoDelete: tododelete.Run,
	}
}
