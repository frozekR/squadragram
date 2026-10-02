package mediametadata_test

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"squadraton-backend/internal/model"
	characterrepository "squadraton-backend/internal/repository/character"
	mediarepository "squadraton-backend/internal/repository/media_metadata"
	skillrepository "squadraton-backend/internal/repository/skill"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// Set TEST_DATABASE_URL to run against PostgreSQL.
// The migration and records live in an isolated schema and are rolled back together.
func TestPostgresMediaOwnershipAndHydration(t *testing.T) {
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("TEST_DATABASE_URL is not set")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	conn, err := pgx.Connect(ctx, dsn)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close(context.Background())
	var version string
	if err := conn.QueryRow(ctx, "SELECT version()").Scan(&version); err != nil {
		t.Fatal(err)
	}
	t.Log(version)
	tx, err := conn.Begin(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback(context.Background())
	schema := pgx.Identifier{"media_test_" + uuid.New().String()}.Sanitize()
	if _, err := tx.Exec(ctx, "CREATE SCHEMA "+schema); err != nil {
		t.Fatal(err)
	}
	if _, err := tx.Exec(ctx, "SET LOCAL search_path TO "+schema); err != nil {
		t.Fatal(err)
	}
	migration, err := os.ReadFile("../../../migrations/000001_init.up.sql")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := tx.Exec(ctx, string(migration)); err != nil {
		t.Fatal(err)
	}
	orderMigration, err := os.ReadFile("../../../migrations/000002_skill_order.up.sql")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := tx.Exec(ctx, string(orderMigration)); err != nil {
		t.Fatal(err)
	}
	characters, skills, media := characterrepository.NewRepository(tx), skillrepository.NewRepository(tx), mediarepository.NewRepository(tx)
	character, err := characters.CreateCharacter(ctx, model.CreateCharacterDTO{Name: "media-test-" + uuid.NewString(), Description: "integration", Role: model.CharacterRoleTank})
	if err != nil {
		t.Fatal(err)
	}
	skill, err := skills.CreateSkill(ctx, model.CreateSkillDTO{Name: "media-test-" + uuid.NewString(), Description: "integration", SkillType: model.SkillTypeSkill, CharacterUUID: character.UUID})
	if err != nil {
		t.Fatal(err)
	}
	for _, item := range []model.MediaMetadata{
		{OwnerType: model.OwnerTypeCharacter, OwnerUUID: character.UUID, MediaType: model.MediaTypeIcon, ObjectKey: "test/character/icon"},
		{OwnerType: model.OwnerTypeCharacter, OwnerUUID: character.UUID, MediaType: model.MediaTypeRender, ObjectKey: "test/character/render"},
		{OwnerType: model.OwnerTypeSkill, OwnerUUID: skill.UUID, MediaType: model.MediaTypeDemo, ObjectKey: "test/skill/demo"},
	} {
		exists, err := media.OwnerExists(ctx, item.OwnerType, item.OwnerUUID)
		if err != nil || !exists {
			t.Fatalf("owner lookup: %v, %v", exists, err)
		}
		created, err := media.CreateMediaMetadata(ctx, item)
		if err != nil {
			t.Fatal(err)
		}
		byID, err := media.GetMediaMetadataByID(ctx, created.ID)
		if err != nil || byID != created {
			t.Fatalf("read by ID: %+v, %v", byID, err)
		}
	}
	character, err = characters.GetCharacterByUUID(ctx, character.UUID)
	if err != nil {
		t.Fatal(err)
	}
	if len(character.MediaMetadata) != 2 {
		t.Fatalf("character media: %+v", character.MediaMetadata)
	}
	skill, err = skills.GetSkillByUUID(ctx, skill.UUID)
	if err != nil {
		t.Fatal(err)
	}
	if len(skill.MediaMetadata) != 1 || skill.MediaMetadata[0].MediaType != model.MediaTypeDemo {
		t.Fatalf("skill media: %+v", skill.MediaMetadata)
	}
	allCharacters, err := characters.GetCharacters(ctx)
	if err != nil {
		t.Fatal(err)
	}
	found := false
	for _, item := range allCharacters {
		if item.UUID == character.UUID {
			found = true
			if len(item.MediaMetadata) != 2 {
				t.Fatal("character list did not load media")
			}
		}
	}
	if !found {
		t.Fatal("character missing from list")
	}
	allSkills, err := skills.GetSkillsByCharacterUUID(ctx, character.UUID)
	if err != nil || len(allSkills) != 1 || len(allSkills[0].MediaMetadata) != 1 {
		t.Fatalf("nested skills media: %+v, %v", allSkills, err)
	}
	wrongType, err := media.GetMediaMetadataByOwner(ctx, model.OwnerTypeSkill, character.UUID)
	if err != nil || len(wrongType) != 0 {
		t.Fatalf("media leaked between owner types: %+v, %v", wrongType, err)
	}
	exists, err := media.OwnerExists(ctx, model.OwnerTypeCharacter, uuid.New())
	if err != nil || exists {
		t.Fatalf("nonexistent owner: %v, %v", exists, err)
	}
	_, err = media.OwnerExists(ctx, model.OwnerType("UNKNOWN"), uuid.New())
	if !errors.Is(err, mediarepository.ErrUnsupportedOwner) {
		t.Fatalf("unsupported owner: %v", err)
	}
	// A savepoint keeps the outer transaction usable after testing the unique constraint.
	duplicateTx, err := tx.Begin(ctx)
	if err != nil {
		t.Fatal(err)
	}
	_, err = mediarepository.NewRepository(duplicateTx).CreateMediaMetadata(ctx, model.MediaMetadata{
		OwnerType: model.OwnerTypeCharacter, OwnerUUID: character.UUID, MediaType: model.MediaTypeIcon, ObjectKey: "test/duplicate",
	})
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) || pgErr.Code != "23505" {
		t.Fatalf("expected unique violation, got %v", err)
	}
	if err := duplicateTx.Rollback(ctx); err != nil {
		t.Fatal(err)
	}
	updated, err := characters.UpdateCharacter(ctx, character.UUID, model.CreateCharacterDTO{Name: "edited-" + uuid.NewString(), Description: "edited", Role: model.CharacterRoleDamage})
	if err != nil || updated.UUID != character.UUID || updated.ID != character.ID || updated.Description != "edited" || len(updated.MediaMetadata) != 2 {
		t.Fatalf("update character: %+v, %v", updated, err)
	}
	other, err := characters.CreateCharacter(ctx, model.CreateCharacterDTO{Name: "other-" + uuid.NewString(), Description: "other", Role: model.CharacterRoleTank})
	if err != nil {
		t.Fatal(err)
	}
	moved, err := skills.UpdateSkill(ctx, skill.UUID, model.CreateSkillDTO{Name: skill.Name, Description: "edited skill", SkillType: model.SkillTypeSkill, CharacterUUID: other.UUID})
	if err != nil || moved.UUID != skill.UUID || moved.CharacterID == nil || *moved.CharacterID != other.ID || len(moved.MediaMetadata) != 1 {
		t.Fatalf("move skill preserving media: %+v, %v", moved, err)
	}
	oldSkills, err := skills.GetSkillsByCharacterUUID(ctx, character.UUID)
	if err != nil || len(oldSkills) != 0 {
		t.Fatalf("old owner retained skill: %+v, %v", oldSkills, err)
	}
	newSkills, err := skills.GetSkillsByCharacterUUID(ctx, other.UUID)
	if err != nil || len(newSkills) != 1 {
		t.Fatalf("new owner lost skill: %+v, %v", newSkills, err)
	}
	m := character.MediaMetadata[0]
	replaced, oldKey, err := media.ReplaceMediaMetadata(ctx, m.ID, "test/replacement")
	if err != nil || oldKey != m.ObjectKey || replaced.ID != m.ID || replaced.OwnerUUID != m.OwnerUUID || replaced.ObjectKey != "test/replacement" {
		t.Fatalf("replace metadata: %+v, %s, %v", replaced, oldKey, err)
	}
	_, oldKey, err = media.ReplaceMediaMetadata(ctx, m.ID, "test/replacement-again")
	if err != nil || oldKey != "test/replacement" {
		t.Fatalf("second replacement cleanup key: %s, %v", oldKey, err)
	}
	deleted, err := media.DeleteMediaMetadata(ctx, m.ID)
	if err != nil || deleted.ObjectKey != "test/replacement-again" {
		t.Fatalf("delete metadata: %+v, %v", deleted, err)
	}
	_, err = media.GetMediaMetadataByID(ctx, m.ID)
	if !errors.Is(err, pgx.ErrNoRows) {
		t.Fatalf("deleted metadata still available: %v", err)
	}
	_, err = skills.UpdateSkill(ctx, skill.UUID, model.CreateSkillDTO{Name: skill.Name, Description: "missing owner", SkillType: model.SkillTypeSkill, CharacterUUID: uuid.New()})
	if !errors.Is(err, pgx.ErrNoRows) {
		t.Fatalf("missing skill owner: %v", err)
	}
	order := 20
	moved, err = skills.UpdateSkill(ctx, skill.UUID, model.CreateSkillDTO{Name: skill.Name, Description: "ordered", SkillType: model.SkillTypeSkill, CharacterUUID: other.UUID, SortOrder: &order})
	if err != nil || moved.SortOrder != 20 {
		t.Fatalf("save skill order: %+v, %v", moved, err)
	}
	// Older clients omit sort_order: an unrelated edit must retain the saved position.
	moved, err = skills.UpdateSkill(ctx, skill.UUID, model.CreateSkillDTO{Name: skill.Name, Description: "keep order", SkillType: model.SkillTypeSkill, CharacterUUID: other.UUID})
	if err != nil || moved.SortOrder != 20 {
		t.Fatalf("omitted order was reset: %+v, %v", moved, err)
	}
	order = 5
	earlier, err := skills.CreateSkill(ctx, model.CreateSkillDTO{Name: "earlier-" + uuid.NewString(), Description: "earlier", SkillType: model.SkillTypeSkill, CharacterUUID: other.UUID, SortOrder: &order})
	if err != nil {
		t.Fatal(err)
	}
	tied, err := skills.CreateSkill(ctx, model.CreateSkillDTO{Name: "tied-" + uuid.NewString(), Description: "same order", SkillType: model.SkillTypeSkill, CharacterUUID: other.UUID, SortOrder: &order})
	if err != nil {
		t.Fatal(err)
	}
	ordered, err := skills.GetSkillsByCharacterUUID(ctx, other.UUID)
	if err != nil || len(ordered) != 3 || ordered[0].UUID != earlier.UUID || ordered[1].UUID != tied.UUID || ordered[2].UUID != moved.UUID {
		t.Fatalf("ordered list with stable ties: %+v, %v", ordered, err)
	}
}
