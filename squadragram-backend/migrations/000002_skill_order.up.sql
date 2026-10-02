ALTER TABLE skills ADD COLUMN IF NOT EXISTS sort_order INTEGER NOT NULL DEFAULT 0
    CHECK (sort_order >= 0);
CREATE INDEX IF NOT EXISTS skills_character_order_idx ON skills (char_id, sort_order, id);
