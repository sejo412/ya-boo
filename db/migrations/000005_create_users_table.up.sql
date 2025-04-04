CREATE TABLE IF NOT EXISTS users(
    id BIGINT PRIMARY KEY UNIQUE NOT NULL,
    first_name VARCHAR NOT NULL,
    last_name VARCHAR,
    username VARCHAR,
    role INT REFERENCES roles (id) ON DELETE SET DEFAULT DEFAULT 0,
    created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    llm INT REFERENCES llm (id) ON DELETE SET DEFAULT DEFAULT 0
);