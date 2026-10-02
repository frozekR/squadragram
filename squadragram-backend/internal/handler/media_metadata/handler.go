package mediametadata

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"squadraton-backend/internal/handler/httpresponse"
	"squadraton-backend/internal/model"
	mediarepository "squadraton-backend/internal/repository/media_metadata"
	storage "squadraton-backend/pkg/minio"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
)

const MaxFileSize int64 = 32 << 20

type Repository interface {
	OwnerExists(context.Context, model.OwnerType, uuid.UUID) (bool, error)
	CreateMediaMetadata(context.Context, model.MediaMetadata) (model.MediaMetadata, error)
	ReplaceMediaMetadata(context.Context, int, string) (model.MediaMetadata, string, error)
	DeleteMediaMetadata(context.Context, int) (model.MediaMetadata, error)
	GetMediaMetadataByID(context.Context, int) (model.MediaMetadata, error)
	GetMediaMetadataByOwner(context.Context, model.OwnerType, uuid.UUID) ([]model.MediaMetadata, error)
}

type Storage interface {
	Upload(context.Context, string, io.Reader, int64, string) error
	Remove(context.Context, string) error
	Open(context.Context, string) (io.ReadCloser, string, error)
}

type Handler struct {
	repository Repository
	storage    Storage
}

func NewHandler(repository Repository, storage Storage) *Handler {
	return &Handler{repository: repository, storage: storage}
}

func (h *Handler) RegisterRoutes(router fiber.Router) {
	router.Post("/media-metadata", h.Upload)
	router.Put("/media-metadata/:id", h.Replace)
	router.Delete("/media-metadata/:id", h.Delete)
	router.Get("/media-metadata", h.GetByOwner)
	router.Get("/media-metadata/:id/file", h.GetFile)
}

func (h *Handler) checkOwner(c fiber.Ctx, ownerType model.OwnerType, u uuid.UUID) (bool, error) {
	exists, err := h.repository.OwnerExists(c.Context(), ownerType, u)
	if errors.Is(err, mediarepository.ErrUnsupportedOwner) {
		return false, httpresponse.Error(c, 400, "unsupported owner_type")
	}
	if err != nil {
		return false, httpresponse.DatabaseError(c, err, "owner not found")
	}
	if !exists {
		return false, httpresponse.Error(c, 404, "owner not found")
	}
	return true, nil
}

func (h *Handler) Upload(c fiber.Ctx) error {
	return h.upload(c, nil)
}

func (h *Handler) Replace(c fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil || id <= 0 {
		return httpresponse.Error(c, 400, "invalid media ID")
	}
	metadata, err := h.repository.GetMediaMetadataByID(c.Context(), id)
	if err != nil {
		return httpresponse.DatabaseError(c, err, "media not found")
	}
	return h.upload(c, &metadata)
}

func (h *Handler) upload(c fiber.Ctx, replacing *model.MediaMetadata) error {
	u, err := uuid.Parse(c.FormValue("owner_uuid"))
	if err != nil || u == uuid.Nil {
		return httpresponse.Error(c, 400, "invalid owner_uuid")
	}
	ownerType := model.OwnerType(c.FormValue("owner_type"))
	mediaType := model.MediaType(c.FormValue("media_type"))
	if replacing != nil && (replacing.OwnerType != ownerType || replacing.OwnerUUID != u || replacing.MediaType != mediaType) {
		return httpresponse.Error(c, 400, "replacement must preserve media owner and type")
	}
	if mediaType != model.MediaTypeIcon && mediaType != model.MediaTypeRender && mediaType != model.MediaTypeDemo {
		return httpresponse.Error(c, 400, "invalid media_type")
	}
	file, err := c.FormFile("file")
	if err != nil {
		return httpresponse.Error(c, 400, "file is required (multipart/form-data)")
	}
	if file.Size == 0 {
		return httpresponse.Error(c, 400, "file is empty")
	}
	if file.Size > MaxFileSize {
		return httpresponse.Error(c, 413, "file exceeds 32 MiB")
	}
	reader, err := file.Open()
	if err != nil {
		return httpresponse.Error(c, 500, "cannot read uploaded file")
	}
	defer reader.Close()
	// Тип определяется по содержимому файла, а не по имени или заголовку клиента.
	var header [512]byte
	n, err := reader.Read(header[:])
	if err != nil && err != io.EOF {
		return httpresponse.Error(c, 500, "cannot read uploaded file")
	}
	contentType := http.DetectContentType(header[:n])
	if !validContentType(mediaType, contentType) {
		return httpresponse.Error(c, 415, "ICON/RENDER require PNG, JPEG, GIF or WebP; DEMO requires MP4 or WebM")
	}
	if _, err := reader.Seek(0, io.SeekStart); err != nil {
		return httpresponse.Error(c, 500, "cannot read uploaded file")
	}
	if ok, err := h.checkOwner(c, ownerType, u); !ok {
		return err
	}
	ctx, cancel := context.WithTimeout(c.Context(), 60*time.Second)
	defer cancel()
	// Случайный суффикс защищает существующий файл от перезаписи при повторной загрузке.
	// UUID владельца остаётся частью ключа для наглядной структуры объектов в MinIO.
	key := fmt.Sprintf("%s/%s/%s/%s", strings.ToLower(string(ownerType)), u,
		strings.ToLower(string(mediaType)), uuid.New())
	if err := h.storage.Upload(ctx, key, reader, file.Size, contentType); err != nil {
		log.Printf("upload media: %v", err)
		return httpresponse.Error(c, 502, "media storage unavailable")
	}
	var metadata model.MediaMetadata
	var previousKey string
	status := fiber.StatusCreated
	if replacing == nil {
		metadata, err = h.repository.CreateMediaMetadata(ctx, model.MediaMetadata{
			OwnerType: ownerType, OwnerUUID: u, MediaType: mediaType, ObjectKey: key,
		})
	} else {
		metadata, previousKey, err = h.repository.ReplaceMediaMetadata(ctx, replacing.ID, key)
		status = fiber.StatusOK
	}
	if err != nil {
		// PostgreSQL и MinIO не имеют общей транзакции. При ошибке записи
		// удаляем только новый объект, сохраняя ранее прикреплённое медиа.
		cleanupCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if cleanupErr := h.storage.Remove(cleanupCtx, key); cleanupErr != nil {
			log.Printf("cleanup media %s: %v", key, cleanupErr)
		}
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return httpresponse.Error(c, 409, "media of this type already exists for owner")
		}
		return httpresponse.DatabaseError(c, err, "owner not found")
	}
	if previousKey != "" {
		h.cleanup(previousKey)
	}
	return c.Status(status).JSON(metadata)
}

func (h *Handler) cleanup(key string) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := h.storage.Remove(ctx, key); err != nil {
		log.Printf("cleanup media %s: %v", key, err)
	}
}

func (h *Handler) Delete(c fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil || id <= 0 {
		return httpresponse.Error(c, 400, "invalid media ID")
	}
	metadata, err := h.repository.DeleteMediaMetadata(c.Context(), id)
	if err != nil {
		return httpresponse.DatabaseError(c, err, "media not found")
	}
	// Unlink first so the API never keeps a reference to an object we already deleted.
	// Storage cleanup is best effort and logged; the metadata deletion has succeeded.
	h.cleanup(metadata.ObjectKey)
	return c.SendStatus(fiber.StatusNoContent)
}

func (h *Handler) GetByOwner(c fiber.Ctx) error {
	u, err := uuid.Parse(c.Query("owner_uuid"))
	if err != nil || u == uuid.Nil {
		return httpresponse.Error(c, 400, "invalid owner_uuid")
	}
	ownerType := model.OwnerType(c.Query("owner_type"))
	if ok, err := h.checkOwner(c, ownerType, u); !ok {
		return err
	}
	metadata, err := h.repository.GetMediaMetadataByOwner(c.Context(), ownerType, u)
	if err != nil {
		return httpresponse.DatabaseError(c, err, "owner not found")
	}
	if metadata == nil {
		metadata = []model.MediaMetadata{}
	}
	return c.JSON(metadata)
}

func (h *Handler) GetFile(c fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil || id <= 0 {
		return httpresponse.Error(c, 400, "invalid media ID")
	}
	metadata, err := h.repository.GetMediaMetadataByID(c.Context(), id)
	if err != nil {
		return httpresponse.DatabaseError(c, err, "media not found")
	}
	reader, contentType, err := h.storage.Open(c.Context(), metadata.ObjectKey)
	if errors.Is(err, storage.ErrObjectNotFound) {
		return httpresponse.Error(c, 404, "media file not found")
	}
	if err != nil {
		log.Printf("open media: %v", err)
		return httpresponse.Error(c, 502, "media storage unavailable")
	}
	c.Set("Content-Type", contentType)
	c.Set("Cache-Control", "no-cache")
	c.Set("X-Content-Type-Options", "nosniff")
	// Fiber закрывает io.ReadCloser после отправки потока, а не при выходе из хэндлера.
	return c.SendStream(reader)
}

func validContentType(mediaType model.MediaType, contentType string) bool {
	if mediaType == model.MediaTypeDemo {
		return contentType == "video/mp4" || contentType == "video/webm"
	}
	return contentType == "image/png" || contentType == "image/jpeg" || contentType == "image/gif" || contentType == "image/webp"
}
