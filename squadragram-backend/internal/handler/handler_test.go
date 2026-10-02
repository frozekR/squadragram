package handler_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	characterhandler "squadraton-backend/internal/handler/character"
	skillhandler "squadraton-backend/internal/handler/skill"
	"squadraton-backend/internal/model"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type characterStub struct {
	characterhandler.Repository
	create func(model.CreateCharacterDTO) (model.Character, error)
	get    func(uuid.UUID) (model.Character, error)
	list   func() ([]model.Character, error)
}

func (s characterStub) CreateCharacter(_ context.Context, dto model.CreateCharacterDTO) (model.Character, error) {
	return s.create(dto)
}

func (s characterStub) GetCharacterByUUID(_ context.Context, u uuid.UUID) (model.Character, error) {
	return s.get(u)
}

func (s characterStub) GetCharacters(context.Context) ([]model.Character, error) {
	return s.list()
}

type skillStub struct {
	skillhandler.Repository
	create      func(model.CreateSkillDTO) (model.Skill, error)
	get         func(uuid.UUID) (model.Skill, error)
	list        func() ([]model.Skill, error)
	byCharacter func(uuid.UUID) ([]model.Skill, error)
}

func (s skillStub) CreateSkill(_ context.Context, dto model.CreateSkillDTO) (model.Skill, error) {
	return s.create(dto)
}

func (s skillStub) GetSkillByUUID(_ context.Context, u uuid.UUID) (model.Skill, error) {
	return s.get(u)
}

func (s skillStub) GetSkills(context.Context) ([]model.Skill, error) {
	return s.list()
}

func (s skillStub) GetSkillsByCharacterUUID(_ context.Context, u uuid.UUID) ([]model.Skill, error) {
	return s.byCharacter(u)
}

func request(t *testing.T, app *fiber.App, method, path, body string, status int) []byte {
	t.Helper()
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}
	response, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()
	data, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatal(err)
	}
	if response.StatusCode != status {
		t.Fatalf("status = %d, want %d; body = %s", response.StatusCode, status, data)
	}
	if !strings.HasPrefix(response.Header.Get("Content-Type"), "application/json") {
		t.Fatalf("expected JSON response, got %q", response.Header.Get("Content-Type"))
	}
	return data
}

func TestCharacterCreateAndRead(t *testing.T) {
	u := uuid.New()
	want := model.Character{ID: 1, UUID: u, Name: "Hero", Description: "Test hero", Role: model.CharacterRoleTank}
	app := fiber.New()
	characterhandler.NewHandler(characterStub{
		create: func(dto model.CreateCharacterDTO) (model.Character, error) {
			if dto.Name != want.Name || dto.Description != want.Description || dto.Role != want.Role {
				t.Fatalf("unexpected DTO: %+v", dto)
			}
			return want, nil
		},
		get: func(got uuid.UUID) (model.Character, error) {
			if got != u {
				t.Fatalf("UUID = %s, want %s", got, u)
			}
			return want, nil
		},
		list: func() ([]model.Character, error) { return []model.Character{want}, nil },
	}).RegisterRoutes(app.Group("/api"))

	for _, tc := range []struct {
		method, path, body string
		status             int
	}{
		{http.MethodPost, "/api/characters", `{"name":"Hero","description":"Test hero","role":"TANK"}`, 201},
		{http.MethodGet, "/api/characters/" + u.String(), "", 200},
	} {
		data := request(t, app, tc.method, tc.path, tc.body, tc.status)
		var got model.Character
		if err := json.Unmarshal(data, &got); err != nil {
			t.Fatal(err)
		}
		if got != want {
			t.Fatalf("character = %+v, want %+v", got, want)
		}
	}
	data := request(t, app, http.MethodGet, "/api/characters", "", 200)
	var got []model.Character
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || got[0] != want {
		t.Fatalf("characters = %+v", got)
	}
}

func TestSkillCreateAndRead(t *testing.T) {
	characterUUID, skillUUID := uuid.New(), uuid.New()
	characterID := 7
	want := model.Skill{ID: 2, UUID: skillUUID, Name: "Shield", Description: "Protect hero", SkillType: model.SkillTypeSkill, CharacterID: &characterID}
	app := fiber.New()
	skillhandler.NewHandler(skillStub{
		create: func(dto model.CreateSkillDTO) (model.Skill, error) {
			if dto.CharacterUUID != characterUUID || dto.Name != want.Name || dto.Description != want.Description || dto.SkillType != want.SkillType {
				t.Fatalf("unexpected DTO: %+v", dto)
			}
			return want, nil
		},
		get: func(u uuid.UUID) (model.Skill, error) {
			if u != skillUUID {
				t.Fatalf("unexpected skill UUID: %s", u)
			}
			return want, nil
		},
		list: func() ([]model.Skill, error) { return []model.Skill{want}, nil },
		byCharacter: func(u uuid.UUID) ([]model.Skill, error) {
			if u != characterUUID {
				t.Fatalf("unexpected character UUID: %s", u)
			}
			return []model.Skill{want}, nil
		},
	}).RegisterRoutes(app.Group("/api"))
	for _, tc := range []struct {
		method, path, body string
		status             int
	}{
		{http.MethodPost, "/api/skills", fmt.Sprintf(`{"name":"Shield","description":"Protect hero","type":"SKILL","character_uuid":"%s"}`, characterUUID), 201},
		{http.MethodGet, "/api/skills/" + skillUUID.String(), "", 200},
	} {
		data := request(t, app, tc.method, tc.path, tc.body, tc.status)
		var got model.Skill
		if err := json.Unmarshal(data, &got); err != nil {
			t.Fatal(err)
		}
		assertSkill(t, got, want)
	}
	for _, path := range []string{"/api/skills", "/api/characters/" + characterUUID.String() + "/skills"} {
		data := request(t, app, http.MethodGet, path, "", 200)
		var got []model.Skill
		if err := json.Unmarshal(data, &got); err != nil {
			t.Fatal(err)
		}
		if len(got) != 1 {
			t.Fatalf("skills = %+v", got)
		}
		assertSkill(t, got[0], want)
	}
}

func assertSkill(t *testing.T, got, want model.Skill) {
	t.Helper()
	if got.CharacterID == nil || *got.CharacterID != *want.CharacterID {
		t.Fatalf("unexpected character link: %+v", got)
	}
	got.CharacterID = want.CharacterID
	if got != want {
		t.Fatalf("skill = %+v, want %+v", got, want)
	}
}

func TestInvalidRequestsDoNotCallRepository(t *testing.T) {
	app := fiber.New()
	characterhandler.NewHandler(characterStub{}).RegisterRoutes(app.Group("/api"))
	skillhandler.NewHandler(skillStub{}).RegisterRoutes(app.Group("/api"))
	u := uuid.New().String()
	for _, tc := range []struct{ name, method, path, body string }{
		{"character UUID", "GET", "/api/characters/bad", ""},
		{"skill UUID", "GET", "/api/skills/bad", ""},
		{"nested UUID", "GET", "/api/characters/bad/skills", ""},
		{"malformed JSON", "POST", "/api/characters", `{"name":`},
		{"missing name", "POST", "/api/characters", `{"description":"x","role":"TANK"}`},
		{"blank name", "POST", "/api/characters", `{"name":"  ","description":"x","role":"TANK"}`},
		{"long name", "POST", "/api/characters", fmt.Sprintf(`{"name":"%s","description":"x","role":"TANK"}`, strings.Repeat("я", 97))},
		{"invalid role", "POST", "/api/characters", `{"name":"x","description":"x","role":"UNKNOWN"}`},
		{"missing character", "POST", "/api/skills", `{"name":"x","description":"x","type":"SKILL"}`},
		{"malformed character UUID", "POST", "/api/skills", `{"name":"x","description":"x","type":"SKILL","character_uuid":"bad"}`},
		{"invalid skill type", "POST", "/api/skills", fmt.Sprintf(`{"name":"x","description":"x","type":"UNKNOWN","character_uuid":"%s"}`, u)},
		{"missing skill description", "POST", "/api/skills", fmt.Sprintf(`{"name":"x","type":"SKILL","character_uuid":"%s"}`, u)},
	} {
		t.Run(tc.name, func(t *testing.T) { request(t, app, tc.method, tc.path, tc.body, 400) })
	}
}

func TestRepositoryErrors(t *testing.T) {
	u := uuid.New().String()
	for _, tc := range []struct {
		name   string
		err    error
		status int
	}{
		{"not found", fmt.Errorf("wrapped: %w", pgx.ErrNoRows), 404},
		{"duplicate name", fmt.Errorf("wrapped: %w", &pgconn.PgError{Code: "23505"}), 409},
		{"missing character FK", &pgconn.PgError{Code: "23503"}, 404},
		{"database failure", errors.New("secret database details"), 500},
	} {
		t.Run(tc.name, func(t *testing.T) {
			app := fiber.New()
			characterhandler.NewHandler(characterStub{
				get: func(uuid.UUID) (model.Character, error) { return model.Character{}, tc.err },
			}).RegisterRoutes(app.Group("/api"))
			skillhandler.NewHandler(skillStub{
				create: func(model.CreateSkillDTO) (model.Skill, error) { return model.Skill{}, tc.err },
			}).RegisterRoutes(app.Group("/api"))
			for _, data := range [][]byte{
				request(t, app, "GET", "/api/characters/"+u, "", tc.status),
				request(t, app, "POST", "/api/skills", fmt.Sprintf(`{"name":"x","description":"x","type":"SKILL","character_uuid":"%s"}`, u), tc.status),
			} {
				var body map[string]string
				if err := json.Unmarshal(data, &body); err != nil {
					t.Fatal(err)
				}
				if body["error"] == "" || strings.Contains(string(data), "secret") {
					t.Fatalf("unexpected error response: %s", data)
				}
			}
		})
	}
}

func TestEmptyListsAreArrays(t *testing.T) {
	app := fiber.New()
	characterhandler.NewHandler(characterStub{list: func() ([]model.Character, error) { return nil, nil }}).RegisterRoutes(app.Group("/api"))
	skillhandler.NewHandler(skillStub{
		list:        func() ([]model.Skill, error) { return nil, nil },
		byCharacter: func(uuid.UUID) ([]model.Skill, error) { return nil, nil },
	}).RegisterRoutes(app.Group("/api"))
	for _, path := range []string{"/api/characters", "/api/skills", "/api/characters/" + uuid.New().String() + "/skills"} {
		if data := request(t, app, "GET", path, "", 200); string(data) != "[]" {
			t.Fatalf("expected [], got %s", data)
		}
	}
}
