CREATE TABLE IF NOT EXISTS texts
(
    id       INTEGER PRIMARY KEY,
    text     TEXT    NOT NULL,
    info     TEXT    NOT NULL DEFAULT '',
    user_uid INTEGER NOT NULL
);
