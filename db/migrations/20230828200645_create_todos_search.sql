CREATE VIRTUAL TABLE todos_search USING fts5(id UNINDEXED, task, tokenize="trigram");

INSERT INTO todos_search(id, task) SELECT id, task FROM todos;

-- Trigger to handle INSERT
CREATE TRIGGER trg_insert_todos AFTER INSERT ON todos
BEGIN
    INSERT INTO todos_search(id, task) VALUES (new.id, new.task);
END;

-- Trigger to handle UPDATE
CREATE TRIGGER trg_update_todos AFTER UPDATE ON todos
BEGIN
    UPDATE todos_search SET task = new.task WHERE id = old.id;
END;

-- Trigger to handle DELETE
CREATE TRIGGER trg_delete_todos AFTER DELETE ON todos
BEGIN
    DELETE FROM todos_search WHERE id = old.id;
END;

PRAGMA USER_VERSION = 2;
