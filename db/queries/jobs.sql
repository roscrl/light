-- name: ScheduleNewJob :exec
INSERT INTO jobs (id, name, run_at, arguments)
VALUES (?, ?, ?, ?);

-- name: GetOverduePendingJobsFromTime :many
SELECT id, name, arguments, run_at
FROM jobs
WHERE run_at <= @from
  AND status = 'pending';

-- name: SetJobStatusToRunning :exec
UPDATE jobs
SET status = 'running'
WHERE id = ?;

-- name: SetFailedJob :exec
UPDATE jobs
SET finished_at    = strftime('%s', 'now'),
    failed_message = ?,
    status         = 'failed'
WHERE id = ?;

-- name: SetSuccessfulJob :exec
UPDATE jobs
SET finished_at = strftime('%s', 'now'),
    status      = 'success'
WHERE id = ?;

-- name: GetJobs :many
SELECT *
FROM jobs
ORDER BY run_at,
         CASE status
             WHEN 'running' THEN 1
             WHEN 'pending' THEN 2
             WHEN 'failed' THEN 3
             WHEN 'success' THEN 4
             ELSE 5
             END
LIMIT ?;

-- name: GetJobsAfterEqualRunAtDate :many
SELECT *
FROM jobs
WHERE run_at >= @from
ORDER BY run_at
LIMIT ?;

-- name: GetJobByID :one
SELECT *
FROM jobs
WHERE id = ?;