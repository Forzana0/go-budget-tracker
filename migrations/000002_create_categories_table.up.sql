CREATE TABLE IF NOT EXISTS categories (
    id      SERIAL PRIMARY KEY,
    name    VARCHAR(30) NOT NULL,
    user_id INTEGER REFERENCES users(id)
);
