CREATE TABLE IF NOT EXISTS llm(
    id serial PRIMARY KEY,
    name VARCHAR(100) UNIQUE NOT NULL,
    endpoint VARCHAR,
    token VARCHAR,
    description VARCHAR(200)

);