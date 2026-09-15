package model

import "github.com/google/uuid"

type MediaMetadata struct {
	ID        int       `json:"id"`
	OwnerType OwnerType `json:"owner_type"`
	OwnerUUID uuid.UUID `json:"owner_uuid"`
	MediaType MediaType `json:"media_type"`
	ObjectKey string    `json:"object_key"`
}
