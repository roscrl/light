package jobs

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/roscrl/light/core/helpers/ulid"
	"github.com/roscrl/light/db/sqlc"
)

func Enqueue(ctx context.Context, name JobName, args map[string]any, qry *sqlc.Queries) (string, error) {
	return Schedule(ctx, name, args, time.Now(), qry)
}

func Schedule(ctx context.Context, name JobName, args map[string]any, t time.Time, qry *sqlc.Queries) (string, error) {
	var argsJSON []byte

	if args != nil {
		var err error
		argsJSON, err = json.Marshal(args)

		if err != nil {
			return "", fmt.Errorf("marshalling args: %w", err)
		}
	}

	id := ulid.NewString()

	scheduleNewJobParams := sqlc.ScheduleNewJobParams{
		ID:        id,
		Name:      string(name),
		RunAt:     t.Unix(),
		Arguments: string(argsJSON),
	}

	if err := qry.ScheduleNewJob(ctx, scheduleNewJobParams); err != nil {
		return "", fmt.Errorf("scheduling job: %w", err)
	}

	return id, nil
}
