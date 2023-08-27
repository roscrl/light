-- name: CreateTodo :execresult
INSERT INTO todos (id, task, status)
VALUES (?, ?, ?);

-- name: GetTodos :many
SELECT id, task, status, created_at
FROM todos
ORDER BY created_at DESC;

-- name: GetTodo :one
SELECT id, task, status, created_at
FROM todos
WHERE id = ?;

-- name: UpdateTodo :one
UPDATE todos
SET task   = ?,
    status = ?
WHERE id = ?
RETURNING id, task, status, created_at;

-- name: SetTodoStatus :execresult
UPDATE todos
SET status = ?
WHERE id = ?