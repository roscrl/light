CREATE VIRTUAL TABLE todos_search USING fts5(id UNINDEXED, task, tokenize="trigram");

INSERT INTO todos_search(id, task) SELECT id, task FROM todos;

-- trigger_<table name>_on_<operation>[_<purpose or target table>]?

-- Trigger to handle INSERT
CREATE TRIGGER trigger_todos_on_insert_search AFTER INSERT ON todos
BEGIN
    INSERT INTO todos_search(id, task) VALUES (new.id, new.task);
END;

-- Trigger to handle UPDATE
CREATE TRIGGER trigger_todos_on_update_search AFTER UPDATE ON todos
BEGIN
    UPDATE todos_search SET task = new.task WHERE id = old.id;
END;

-- Trigger to handle DELETE
CREATE TRIGGER trigger_todos_on_delete_search AFTER DELETE ON todos
BEGIN
    DELETE FROM todos_search WHERE id = old.id;
END;

PRAGMA USER_VERSION = 4;
