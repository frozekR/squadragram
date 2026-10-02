package mediametadata

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
	"strings"
	"testing"

	"squadraton-backend/internal/model"
	mediarepository "squadraton-backend/internal/repository/media_metadata"
	storage "squadraton-backend/pkg/minio"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type repositoryStub struct {
	Repository
	exists    bool
	ownerErr  error
	createErr error
	metadata  model.MediaMetadata
}

func (r *repositoryStub) OwnerExists(context.Context, model.OwnerType, uuid.UUID) (bool, error) {
	return r.exists, r.ownerErr
}

func (r *repositoryStub) CreateMediaMetadata(_ context.Context, m model.MediaMetadata) (model.MediaMetadata, error) {
	if r.createErr != nil {
		return model.MediaMetadata{}, r.createErr
	}
	m.ID = 17
	r.metadata = m
	return m, nil
}

func (r *repositoryStub) GetMediaMetadataByID(_ context.Context, id int) (model.MediaMetadata, error) {
	if id != r.metadata.ID {
		return model.MediaMetadata{}, pgx.ErrNoRows
	}
	return r.metadata, nil
}

func (r *repositoryStub) ReplaceMediaMetadata(_ context.Context, id int, key string) (model.MediaMetadata, string, error) {
	if r.createErr != nil {
		return model.MediaMetadata{}, "", r.createErr
	}
	if id != r.metadata.ID {
		return model.MediaMetadata{}, "", pgx.ErrNoRows
	}
	old := r.metadata.ObjectKey
	r.metadata.ObjectKey = key
	return r.metadata, old, nil
}

func (r *repositoryStub) DeleteMediaMetadata(_ context.Context, id int) (model.MediaMetadata, error) {
	if id != r.metadata.ID {
		return model.MediaMetadata{}, pgx.ErrNoRows
	}
	m := r.metadata
	r.metadata = model.MediaMetadata{}
	return m, nil
}

func (r *repositoryStub) GetMediaMetadataByOwner(_ context.Context, ownerType model.OwnerType, u uuid.UUID) ([]model.MediaMetadata, error) {
	if r.metadata.OwnerType == ownerType && r.metadata.OwnerUUID == u {
		return []model.MediaMetadata{r.metadata}, nil
	}
	return nil, nil
}

type storageStub struct {
	Storage
	key         string
	removed     string
	data        []byte
	contentType string
	uploadErr   error
	openErr     error
}

func (s *storageStub) Upload(_ context.Context, key string, reader io.Reader, _ int64, contentType string) error {
	s.key, s.contentType = key, contentType
	s.data, _ = io.ReadAll(reader)
	return s.uploadErr
}

func (s *storageStub) Remove(_ context.Context, key string) error { s.removed = key; return nil }

func (s *storageStub) Open(_ context.Context, key string) (io.ReadCloser, string, error) {
	if s.openErr != nil {
		return nil, "", s.openErr
	}
	if key != s.key {
		return nil, "", storage.ErrObjectNotFound
	}
	return io.NopCloser(bytes.NewReader(s.data)), s.contentType, nil
}

func pngFile(t *testing.T) []byte {
	t.Helper()
	var data bytes.Buffer
	if err := png.Encode(&data, image.NewRGBA(image.Rect(0, 0, 2, 2))); err != nil {
		t.Fatal(err)
	}
	return data.Bytes()
}

func uploadRequest(t *testing.T, fields map[string]string, file []byte) *http.Request {
	t.Helper()
	var data bytes.Buffer
	writer := multipart.NewWriter(&data)
	for k, v := range fields {
		if err := writer.WriteField(k, v); err != nil {
			t.Fatal(err)
		}
	}
	if file != nil {
		part, err := writer.CreateFormFile("file", "icon.png")
		if err != nil {
			t.Fatal(err)
		}
		if _, err := part.Write(file); err != nil {
			t.Fatal(err)
		}
	}
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}
	req := httptest.NewRequest("POST", "/api/media-metadata", &data)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req
}

func response(t *testing.T, app *fiber.App, req *http.Request, status int) []byte {
	t.Helper()
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != status {
		t.Fatalf("status %d, want %d: %s", resp.StatusCode, status, data)
	}
	return data
}

func TestUploadMetadataAndFile(t *testing.T) {
	u := uuid.New()
	repo, store := &repositoryStub{exists: true}, &storageStub{}
	app := fiber.New()
	NewHandler(repo, store).RegisterRoutes(app.Group("/api"))
	file := pngFile(t)
	data := response(t, app, uploadRequest(t, map[string]string{
		"owner_type": "CHARACTER", "owner_uuid": u.String(), "media_type": "ICON",
	}, file), 201)
	var got model.MediaMetadata
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatal(err)
	}
	if got.ID != 17 || got.OwnerUUID != u || got.OwnerType != model.OwnerTypeCharacter || got.MediaType != model.MediaTypeIcon || got.ObjectKey != store.key {
		t.Fatalf("unexpected metadata: %+v", got)
	}
	if !strings.HasPrefix(store.key, "character/"+u.String()+"/icon/") || store.contentType != "image/png" || !bytes.Equal(store.data, file) {
		t.Fatal("upload contents or key mismatch")
	}
	data = response(t, app, httptest.NewRequest("GET", "/api/media-metadata?owner_type=CHARACTER&owner_uuid="+u.String(), nil), 200)
	var list []model.MediaMetadata
	if err := json.Unmarshal(data, &list); err != nil {
		t.Fatal(err)
	}
	if len(list) != 1 || list[0] != got {
		t.Fatalf("unexpected list: %s", data)
	}
	data = response(t, app, httptest.NewRequest("GET", "/api/media-metadata/17/file", nil), 200)
	if !bytes.Equal(data, file) {
		t.Fatal("download differs from upload")
	}
}

func TestUploadValidationAndCleanup(t *testing.T) {
	for _, tc := range []struct {
		name                           string
		ownerUUID, mediaType           string
		file                           []byte
		exists                         bool
		ownerErr, createErr, uploadErr error
		status                         int
		cleanup                        bool
	}{
		{name: "invalid UUID", ownerUUID: "bad", mediaType: "ICON", file: pngFile(t), status: 400},
		{name: "invalid type", mediaType: "BAD", file: pngFile(t), status: 400},
		{name: "missing file", mediaType: "ICON", status: 400},
		{name: "empty file", mediaType: "ICON", file: []byte{}, status: 400},
		{name: "text as image", mediaType: "ICON", file: []byte("not an image"), status: 415},
		{name: "image as demo", mediaType: "DEMO", file: pngFile(t), status: 415},
		{name: "owner missing", mediaType: "ICON", file: pngFile(t), status: 404},
		{name: "unsupported owner", mediaType: "ICON", file: pngFile(t), ownerErr: mediarepository.ErrUnsupportedOwner, status: 400},
		{name: "storage failure", mediaType: "ICON", file: pngFile(t), exists: true, uploadErr: errors.New("offline"), status: 502},
		{name: "duplicate media", mediaType: "ICON", file: pngFile(t), exists: true, createErr: &pgconn.PgError{Code: "23505"}, status: 409, cleanup: true},
		{name: "database failure", mediaType: "ICON", file: pngFile(t), exists: true, createErr: errors.New("offline"), status: 500, cleanup: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			u := tc.ownerUUID
			if u == "" {
				u = uuid.New().String()
			}
			repo := &repositoryStub{exists: tc.exists, ownerErr: tc.ownerErr, createErr: tc.createErr}
			store := &storageStub{uploadErr: tc.uploadErr}
			app := fiber.New()
			NewHandler(repo, store).RegisterRoutes(app.Group("/api"))
			response(t, app, uploadRequest(t, map[string]string{"owner_type": "CHARACTER", "owner_uuid": u, "media_type": tc.mediaType}, tc.file), tc.status)
			if tc.cleanup && (store.key == "" || store.removed != store.key) {
				t.Fatal("new object was not removed after database failure")
			}
			if !tc.cleanup && store.removed != "" {
				t.Fatal("unexpected cleanup")
			}
			if tc.status != 500 && tc.status != 409 && tc.status != 502 && store.key != "" {
				t.Fatal("invalid upload reached storage")
			}
		})
	}
}

func TestDemoUploadDetectsMP4(t *testing.T) {
	repo, store := &repositoryStub{exists: true}, &storageStub{}
	app := fiber.New()
	NewHandler(repo, store).RegisterRoutes(app.Group("/api"))
	// A minimal MP4 ftyp box checks content detection independently of the filename.
	file := append([]byte{0, 0, 0, 24}, []byte("ftypmp42\x00\x00\x00\x00mp42isom")...)
	response(t, app, uploadRequest(t, map[string]string{
		"owner_type": "SKILL", "owner_uuid": uuid.NewString(), "media_type": "DEMO",
	}, file), 201)
	if store.contentType != "video/mp4" || repo.metadata.MediaType != model.MediaTypeDemo || repo.metadata.OwnerType != model.OwnerTypeSkill {
		t.Fatalf("unexpected DEMO upload: %s, %+v", store.contentType, repo.metadata)
	}
}

func TestMissingFileResponses(t *testing.T) {
	repo := &repositoryStub{metadata: model.MediaMetadata{ID: 17}}
	store := &storageStub{openErr: storage.ErrObjectNotFound}
	app := fiber.New()
	NewHandler(repo, store).RegisterRoutes(app.Group("/api"))
	response(t, app, httptest.NewRequest("GET", "/api/media-metadata/bad/file", nil), 400)
	response(t, app, httptest.NewRequest("GET", "/api/media-metadata/99/file", nil), 404)
	response(t, app, httptest.NewRequest("GET", "/api/media-metadata/17/file", nil), 404)
}

func TestReplaceAndDeleteMedia(t *testing.T) {
	for _, fail := range []bool{false, true} {
		t.Run(fmt.Sprintf("database_failure_%v", fail), func(t *testing.T) {
			u := uuid.New()
			repo := &repositoryStub{exists: true, metadata: model.MediaMetadata{ID: 17, OwnerType: model.OwnerTypeCharacter, OwnerUUID: u, MediaType: model.MediaTypeIcon, ObjectKey: "old/icon"}}
			store := &storageStub{}
			app := fiber.New()
			NewHandler(repo, store).RegisterRoutes(app.Group("/api"))
			if fail {
				repo.createErr = errors.New("database offline")
			}
			req := uploadRequest(t, map[string]string{"owner_type": "CHARACTER", "owner_uuid": u.String(), "media_type": "ICON"}, pngFile(t))
			req.Method, req.URL.Path = "PUT", "/api/media-metadata/17"
			status := 200
			if fail {
				status = 500
			}
			response(t, app, req, status)
			if fail {
				if repo.metadata.ObjectKey != "old/icon" || store.removed != store.key {
					t.Fatal("failed replacement must preserve old metadata and remove only new object")
				}
				return
			}
			if repo.metadata.ID != 17 || repo.metadata.ObjectKey != store.key || store.removed != "old/icon" {
				t.Fatal("replacement must retain ID and clean old object")
			}
			response(t, app, httptest.NewRequest("DELETE", "/api/media-metadata/17", nil), 204)
			if repo.metadata.ID != 0 || store.removed != store.key {
				t.Fatal("delete must remove metadata and its object")
			}
			response(t, app, httptest.NewRequest("GET", "/api/media-metadata/17/file", nil), 404)
			response(t, app, httptest.NewRequest("DELETE", "/api/media-metadata/17", nil), 404)
		})
	}
}

func TestReplacementCannotChangeOwner(t *testing.T) {
	repo := &repositoryStub{exists: true, metadata: model.MediaMetadata{ID: 17, OwnerType: model.OwnerTypeSkill, OwnerUUID: uuid.New(), MediaType: model.MediaTypeIcon, ObjectKey: "old/icon"}}
	store := &storageStub{}
	app := fiber.New()
	NewHandler(repo, store).RegisterRoutes(app.Group("/api"))
	req := uploadRequest(t, map[string]string{"owner_type": "SKILL", "owner_uuid": uuid.NewString(), "media_type": "ICON"}, pngFile(t))
	req.Method, req.URL.Path = "PUT", "/api/media-metadata/17"
	response(t, app, req, 400)
	if store.key != "" || store.removed != "" || repo.metadata.ObjectKey != "old/icon" {
		t.Fatal("invalid replacement reached storage")
	}
}
