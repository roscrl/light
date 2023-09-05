-- name: CreateTodo :execresult
INSERT INTO todos (id, task, status)
VALUES (?, ?, ?);

-- name: ReadTodos :many
SELECT id, task, status, created_at
FROM todos
ORDER BY created_at DESC;

-- name: ReadTodo :one
SELECT id, task, status, created_at
FROM todos
WHERE id = ?;

-- name: UpdateTodo :one
UPDATE todos
SET task   = ?,
    status = ?
WHERE id = ?
RETURNING id, task, status, created_at;

-- name: SearchTodosOld :many
SELECT t.id, t.task, t.status, t.created_at
FROM todos t
         JOIN todos_search ts ON t.id = ts.id
WHERE ts.task MATCH ?;

-- name: SearchTodos :many
SELECT t.id, t.task, t.status, t.created_at
FROM todos t
         JOIN todos_search ts ON t.id = ts.id
WHERE ts.task LIKE ?
ORDER BY RANK LIMIT 20;