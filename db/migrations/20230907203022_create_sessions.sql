CREATE TABLE sessions
(
    token  TEXT PRIMARY KEY,
    data   BLOB NOT NULL,
    expiry REAL NOT NULL
);

CREATE INDEX index_sessions_expiry ON sessions (expiry);

PRAGMA USER_VERSION = 3;
