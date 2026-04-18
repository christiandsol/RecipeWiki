CREATE TABLE IF NOT EXISTS steps (
    step_id          SERIAL PRIMARY KEY,
    recipe_id   INT NOT NULL REFERENCES recipes(id) ON DELETE CASCADE,
    step_number FLOAT NOT NULL,
    step_text        TEXT NOT NULL,
    UNIQUE (recipe_id, step_number)
);
