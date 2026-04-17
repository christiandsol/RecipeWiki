CREATE TABLE IF NOT EXISTS ingredients (
    id         SERIAL PRIMARY KEY,
    recipe_id  INT NOT NULL REFERENCES recipes(id) ON DELETE CASCADE,
    name       TEXT NOT NULL,
    amount     INT,
    specifier  TEXT,
    current_amount TEXT
);

