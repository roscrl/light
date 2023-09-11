-- name: ScheduleNewJob :exec
INSERT INTO jobs (id, name, run_at, arguments)
VALUES (?, ?, ?, ?);

-- name: GetOverdueJobsFromTime :many
SELECT id, name, arguments, run_at
FROM jobs
WHERE run_at <= @from AND status = 'pending';

-- name: SetJobStatusToRunning :exec
UPDATE jobs
SET status = 'running'
WHERE id = ?;

-- name: SetFailedJob :exec
UPDATE jobs
SET failed_at = strftime('%s', 'now'), failed_message = ?, status = 'failed'
WHERE id = ?;

-- name: SetSuccessfulJob :exec
UPDATE jobs
SET completed_at = strftime('%s', 'now'), status = 'success'
WHERE id = ?;
