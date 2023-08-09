-- name: CreateTodo :execresult
INSERT INTO todos (id, task, status)
VALUES (?, ?, ?);

-- name: GetTodos :many
SELECT id, task, status, created_at
FROM todos
ORDER BY created_at DESC;

-- name: SetTodoStatus :execresult
UPDATE todos
SET status = ?
WHERE id = ?