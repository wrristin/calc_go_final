CREATE TABLE IF NOT EXISTS users (
    id TEXT PRIMARY KEY,
    login TEXT UNIQUE,
    password TEXT
);

CREATE TABLE IF NOT EXISTS expressions (
    id TEXT PRIMARY KEY,
    user_id TEXT,
    status TEXT,
    result TEXT
);

CREATE TABLE IF NOT EXISTS tasks (
    id TEXT PRIMARY KEY,
    expression_id TEXT,
    arg1 REAL,
    arg2 REAL,
    operation TEXT,
    operation_time INTEGER,
    status TEXT,
    result REAL
);