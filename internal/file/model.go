package file

import (
	"io"
	"time"
)

type FileData struct {
	ID         string    `db:"id"`
	Code       string    `db:"code"`
	Token      string    `db:"token"`
	Name       string    `db:"name"`
	StorageKey string    `db:"storage_key"`
	Size       uint64    `db:"size"`
	CreatedAt  time.Time `db:"created_at"`
	ExpiresAt  time.Time `db:"expires_at"`
}

type RQFileData struct {
	Name   string        `validate:"required,min=1,max=255"`
	TTL    time.Duration `validate:"omitempty"`
	Reader io.Reader     `validate:"required"`
}

type RSFileData struct {
	ID        string    `json:"id"`
	Code      string    `json:"code"`
	Token     *string   `json:"token,omitempty"`
	Name      string    `json:"name"`
	Size      uint64    `json:"size"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

func (fd FileData) MapToResponse(token *string) RSFileData {
	return RSFileData{
		ID:        fd.ID,
		Code:      fd.Code,
		Token:     token,
		Name:      fd.Name,
		Size:      fd.Size,
		CreatedAt: fd.CreatedAt,
		ExpiresAt: fd.ExpiresAt,
	}
}
