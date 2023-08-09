CREATE TABLE todos
(
    id         TEXT PRIMARY KEY,
    task       TEXT    NOT NULL CHECK (LENGTH(task) >= 1 AND LENGTH(task) <= 255),
    status     TEXT    NOT NULL DEFAULT 'pending' CHECK ( status IN ('pending', 'done') ),
    created_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now'))
) STRICT;