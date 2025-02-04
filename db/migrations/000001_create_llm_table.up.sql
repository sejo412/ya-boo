CREATE TABLE IF NOT EXISTS llm(
    id serial PRIMARY KEY,
    name VARCHAR(20) UNIQUE NOT NULL,
    endpoint VARCHAR,
    token VARCHAR,
    description VARCHAR(100)

);