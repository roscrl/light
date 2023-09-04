CREATE TABLE todos
(
    id         TEXT PRIMARY KEY,
    task       TEXT    NOT NULL CHECK (LENGTH(task) >= 1 AND LENGTH(task) <= 255),
    status     TEXT    NOT NULL DEFAULT 'pending' CHECK ( status IN ('pending', 'done') ),
    created_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now'))
) STRICT;
CREATE TABLE IF NOT EXISTS 'todos_search_data'(id INTEGER PRIMARY KEY, block BLOB);
CREATE TABLE IF NOT EXISTS 'todos_search_idx'(segid, term, pgno, PRIMARY KEY(segid, term)) WITHOUT ROWID;
CREATE TABLE IF NOT EXISTS 'todos_search_content'(id INTEGER PRIMARY KEY, c0, c1);
CREATE TABLE IF NOT EXISTS 'todos_search_docsize'(id INTEGER PRIMARY KEY, sz BLOB);
CREATE TABLE IF NOT EXISTS 'todos_search_config'(k PRIMARY KEY, v) WITHOUT ROWID;
CREATE VIRTUAL TABLE todos_search USING fts5(id UNINDEXED, task)
/* todos_search(id,task) */;
CREATE TRIGGER trg_insert_todos AFTER INSERT ON todos
BEGIN
    INSERT INTO todos_search(id, task) VALUES (new.id, new.task);
END;
CREATE TRIGGER trg_update_todos AFTER UPDATE ON todos
BEGIN
    UPDATE todos_search SET task = new.task WHERE id = old.id;
END;
CREATE TRIGGER trg_delete_todos AFTER DELETE ON todos
BEGIN
    DELETE FROM todos_search WHERE id = old.id;
END;
