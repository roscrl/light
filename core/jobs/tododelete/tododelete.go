package tododelete

import (
	"context"
	"fmt"

	"github.com/roscrl/light/core/jobs/scope"
)

const (
	id = "id"
)

func Args(todoID string) map[string]any {
	return map[string]any{
		id: todoID,
	}
}

func Run(ctx context.Context, sc *scope.Job, args map[string]any) error {
	todoID, ok := args[id].(string)
	if !ok {
		return fmt.Errorf("no todo id given in job args")
	}

	err := sc.Qry.DeleteTodoByID(ctx, todoID)
	if err != nil {
		return fmt.Errorf("deleting todo: %w", err)
	}

	return nil
}
