package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"image/png"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"squadraton-backend/internal/model"
	minio "squadraton-backend/pkg/minio"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// Exercise the same routes as the admin UI with real PostgreSQL and MinIO.
// An isolated transaction rolls back all records; uploaded test objects are removed.
func TestAdminWorkflowIntegration(t *testing.T) {
	dsn, endpoint := os.Getenv("TEST_DATABASE_URL"), os.Getenv("TEST_MINIO_ENDPOINT")
	if dsn == "" || endpoint == "" {
		t.Skip("TEST_DATABASE_URL and TEST_MINIO_ENDPOINT are required")
	}
	t.Setenv("MINIO_ENDPOINT", endpoint)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	conn, err := pgx.Connect(ctx, dsn)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close(context.Background())
	tx, err := conn.Begin(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback(context.Background())
	schema := pgx.Identifier{"admin_test_" + uuid.NewString()}.Sanitize()
	if _, err := tx.Exec(ctx, "CREATE SCHEMA "+schema); err != nil {
		t.Fatal(err)
	}
	if _, err := tx.Exec(ctx, "SET LOCAL search_path TO "+schema); err != nil {
		t.Fatal(err)
	}
	migration, err := os.ReadFile("../migrations/000001_init.up.sql")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := tx.Exec(ctx, string(migration)); err != nil {
		t.Fatal(err)
	}
	orderMigration, err := os.ReadFile("../migrations/000002_skill_order.up.sql")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := tx.Exec(ctx, string(orderMigration)); err != nil {
		t.Fatal(err)
	}
	storage, err := minio.NewMinio()
	if err != nil {
		t.Fatal(err)
	}
	const token = "integration-test-only-admin-token-0123456789"
	app, err := newApp(tx, storage, token)
	if err != nil {
		t.Fatal(err)
	}
	call := func(req *http.Request, status int, target any) []byte {
		t.Helper()
		req.Header.Set("Authorization", "Bearer "+token)
		res, err := app.Test(req)
		if err != nil {
			t.Fatal(err)
		}
		defer res.Body.Close()
		data, err := io.ReadAll(res.Body)
		if err != nil {
			t.Fatal(err)
		}
		if res.StatusCode != status {
			t.Fatalf("%s %s: got %d, want %d: %s", req.Method, req.URL.Path, res.StatusCode, status, data)
		}
		if target != nil {
			if err := json.Unmarshal(data, target); err != nil {
				t.Fatal(err)
			}
		}
		return data
	}
	jsonRequest := func(method, path string, dto any) *http.Request {
		data, err := json.Marshal(dto)
		if err != nil {
			t.Fatal(err)
		}
		req := httptest.NewRequest(method, path, bytes.NewReader(data))
		req.Header.Set("Content-Type", "application/json")
		return req
	}
	call(httptest.NewRequest("POST", "/api/admin/verify", nil), 204, nil)
	heroDTO := model.CreateCharacterDTO{Name: "Admin integration hero", Description: "Before", Role: model.CharacterRoleTank}
	var hero model.Character
	call(jsonRequest("POST", "/api/characters", heroDTO), 201, &hero)
	if hero.UUID == uuid.Nil || hero.ID == 0 {
		t.Fatal("hero identity missing")
	}
	path := "/api/characters/" + hero.UUID.String()
	heroDTO.Description = "After"
	call(jsonRequest("PUT", path, heroDTO), 200, &hero)
	if hero.Description != "After" {
		t.Fatal("hero edit not saved")
	}
	skillDTO := model.CreateSkillDTO{Name: "Admin integration skill", Description: "Before", SkillType: model.SkillTypeSkill, CharacterUUID: hero.UUID}
	var skill model.Skill
	call(jsonRequest("POST", "/api/skills", skillDTO), 201, &skill)
	if skill.CharacterID == nil || *skill.CharacterID != hero.ID {
		t.Fatal("skill not attached to hero")
	}
	skillDTO.Description = "After"
	call(jsonRequest("PUT", "/api/skills/"+skill.UUID.String(), skillDTO), 200, &skill)
	if skill.Description != "After" {
		t.Fatal("skill edit not saved")
	}
	var file bytes.Buffer
	if err := png.Encode(&file, image.NewRGBA(image.Rect(0, 0, 3, 3))); err != nil {
		t.Fatal(err)
	}
	upload := func(method, path string) *http.Request {
		var body bytes.Buffer
		writer := multipart.NewWriter(&body)
		for k, v := range map[string]string{"owner_type": "SKILL", "owner_uuid": skill.UUID.String(), "media_type": "ICON"} {
			if err := writer.WriteField(k, v); err != nil {
				t.Fatal(err)
			}
		}
		part, err := writer.CreateFormFile("file", "test.png")
		if err != nil {
			t.Fatal(err)
		}
		if _, err := part.Write(file.Bytes()); err != nil {
			t.Fatal(err)
		}
		if err := writer.Close(); err != nil {
			t.Fatal(err)
		}
		req := httptest.NewRequest(method, path, &body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		return req
	}
	var media model.MediaMetadata
	call(upload("POST", "/api/media-metadata"), 201, &media)
	cleanup := func(key string) {
		t.Cleanup(func() {
			if err := storage.Remove(context.Background(), key); err != nil {
				t.Errorf("cleanup: %v", err)
			}
		})
	}
	cleanup(media.ObjectKey)
	mediaPath := "/api/media-metadata/" + fmt.Sprint(media.ID)
	oldKey := media.ObjectKey
	call(upload("PUT", mediaPath), 200, &media)
	cleanup(media.ObjectKey)
	if media.ObjectKey == oldKey {
		t.Fatal("replacement did not change object")
	}
	oldReader, _, err := storage.Open(ctx, oldKey)
	if err == nil {
		oldReader.Close()
		t.Fatal("old object was not removed")
	}
	if !errors.Is(err, minio.ErrObjectNotFound) {
		t.Fatalf("old object lookup: %v", err)
	}
	var linked []model.Skill
	call(httptest.NewRequest("GET", path+"/skills", nil), 200, &linked)
	if len(linked) != 1 || len(linked[0].MediaMetadata) != 1 || linked[0].MediaMetadata[0] != media {
		t.Fatalf("skill media hydration: %+v", linked)
	}
	gotFile := call(httptest.NewRequest("GET", mediaPath+"/file", nil), 200, nil)
	if !bytes.Equal(gotFile, file.Bytes()) {
		t.Fatal("download differs from uploaded image")
	}
	call(httptest.NewRequest("DELETE", mediaPath, nil), 204, nil)
	call(httptest.NewRequest("GET", mediaPath+"/file", nil), 404, nil)
	call(httptest.NewRequest("GET", path+"/skills", nil), 200, &linked)
	if len(linked[0].MediaMetadata) != 0 {
		t.Fatal("deleted media still attached")
	}
}
