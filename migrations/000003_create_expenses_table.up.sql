CREATE TABLE IF NOT EXISTS expenses (
    id          SERIAL PRIMARY KEY,
    name        VARCHAR(30)              NOT NULL,
    description VARCHAR(1000),
    price       NUMERIC(10, 2)           NOT NULL,
    created_at  TIMESTAMP WITHOUT TIME ZONE DEFAULT now(),
    user_id     INTEGER REFERENCES users(id),
    category_id INTEGER REFERENCES categories(id)
);
