DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'inventory_level') THEN
        CREATE TYPE inventory_level AS ENUM ('high', 'medium', 'low', 'out');
    END IF;
END$$;


ALTER TABLE ingredients
    ALTER COLUMN current_amount TYPE inventory_level
    USING current_amount::inventory_level;


ALTER TABLE ingredients
    ALTER COLUMN current_amount SET NOT NULL;


ALTER TABLE ingredients
    ALTER COLUMN current_amount SET DEFAULT 'out';
