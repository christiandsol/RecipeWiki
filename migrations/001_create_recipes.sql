CREATE TABLE IF NOT EXISTS recipes (
    id          SERIAL PRIMARY KEY,
    name        TEXT NOT NULL,
    description TEXT,
    steps       TEXT[],
    info        TEXT[]
);

