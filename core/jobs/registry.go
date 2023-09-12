package jobs

import (
	"context"

	"github.com/roscrl/light/core/jobs/scope"
	"github.com/roscrl/light/core/jobs/tododelete"
)

type JobName string

type JobFunc func(ctx context.Context, jobScope *scope.Job, args map[string]any) error

type JobNameToJobFuncRegistry map[JobName]JobFunc

const (
	TodoDelete JobName = "todo_delete"
)

func DefaultRegistry() JobNameToJobFuncRegistry {
	return map[JobName]JobFunc{
		TodoDelete: tododelete.Run,
	}
}
