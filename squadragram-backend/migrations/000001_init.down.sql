DROP TRIGGER IF EXISTS update_characters_updated_at ON characters;
DROP TRIGGER IF EXISTS update_skills_updated_at ON skills;
DROP TRIGGER IF EXISTS update_media_metadata_updated_at ON media_metadata;
DROP FUNCTION IF EXISTS update_updated_at_column();

DROP TABLE IF EXISTS media_metadata;
DROP TABLE IF EXISTS skills;
DROP TABLE IF EXISTS characters;

DROP TYPE IF EXISTS media_type;
DROP TYPE IF EXISTS owner_type;
DROP TYPE IF EXISTS skill_type;
DROP TYPE IF EXISTS character_role;