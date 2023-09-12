-- name: NewTodo :execresult
INSERT INTO todos (id, task, status)
VALUES (?, ?, ?);

-- name: GetAllTodos :many
SELECT id, task, status, created_at
FROM todos
ORDER BY created_at DESC;

-- name: GetTodoByID :one
SELECT id, task, status, created_at
FROM todos
WHERE id = ?;

-- name: UpdateTodoByID :one
UPDATE todos
SET task   = ?,
    status = ?
WHERE id = ?
RETURNING id, task, status, created_at;

-- name: SearchTodosByTask :many
SELECT t.id, t.task, t.status, t.created_at
FROM todos t
         JOIN todos_search ts ON t.id = ts.id
WHERE ts.task LIKE ?
ORDER BY RANK LIMIT 20;

-- name: DeleteTodoByID :exec
DELETE FROM todos
WHERE id = ?;