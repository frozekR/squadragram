CREATE TYPE character_role AS ENUM ('DAMAGE', 'TANK', 'TECHNICAL');
CREATE TYPE skill_type AS ENUM ('PASSIVE', 'RUSH_ATTACK', 'SKILL', 'SUPER_ATTACK', 'TRANSFORMATION');
CREATE TYPE owner_type AS ENUM ('CHARACTER', 'SKILL');
CREATE TYPE media_type AS ENUM ('ICON', 'RENDER', 'DEMO');

CREATE TABLE IF NOT EXISTS characters (
    id BIGSERIAL PRIMARY KEY,
    uuid UUID UNIQUE NOT NULL,
    name VARCHAR(96) UNIQUE NOT NULL,
    description TEXT,
    role character_role NOT NULL,
    
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
);

CREATE TABLE IF NOT EXISTS skills (
    id BIGSERIAL PRIMARY KEY,
    uuid UUID UNIQUE NOT NULL,
    name VARCHAR(96) UNIQUE NOT NULL,
    description TEXT,
    type skill_type NOT NULL,
    char_id BIGINT REFERENCES characters(id),

    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
);

CREATE TABLE IF NOT EXISTS media_metadata (
    id BIGSERIAL PRIMARY KEY,
    owner_type owner_type NOT NULL,
    owner_uuid UUID NOT NULL,
    media_type media_type NOT NULL,
    object_key VARCHAR(512) NOT NULL,

    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now(),
    UNIQUE (owner_type, owner_uuid, media_type)
);

CREATE INDEX idx_skills_char_id ON skills(char_id);
CREATE INDEX idx_media_owner ON media_metadata(owner_type, owner_uuid);

CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_characters_updated_at BEFORE UPDATE ON characters
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_skills_updated_at BEFORE UPDATE ON skills
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_media_metadata_updated_at BEFORE UPDATE ON media_metadata
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();