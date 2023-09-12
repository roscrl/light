CREATE TABLE jobs
(
    id             TEXT PRIMARY KEY,

    name           TEXT    NOT NULL,
    status         TEXT    NOT NULL DEFAULT 'pending' CHECK ( status IN ('pending', 'running', 'success', 'failed') ),
    run_at         INTEGER NOT NULL,
    arguments      TEXT    NOT NULL,
    finished_at    INTEGER,

    failed_message TEXT,

    created_at     INTEGER NOT NULL DEFAULT (strftime('%s', 'now'))
) STRICT;

-- https://docs.gitlab.com/ee/development/database/constraint_naming_convention.html

CREATE INDEX index_jobs_status ON jobs (status);
CREATE INDEX index_jobs_run_at ON jobs (run_at);

PRAGMA USER_VERSION = 2;
