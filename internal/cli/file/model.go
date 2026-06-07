package file

import "time"

type FileOptions struct {
	IsEncrypted bool
	IsPublic    bool
	Token       string
	Password    string
	TTL         string
}

type RSFileData struct {
	ID           string    `json:"id"`
	Code         string    `json:"code"`
	Token        *string   `json:"token,omitempty"`
	Name         string    `json:"name"`
	Size         uint64    `json:"size"`
	Is_Encrypted bool      `json:"is_encrypted"`
	CreatedAt    time.Time `json:"created_at"`
	ExpiresAt    time.Time `json:"expires_at"`
}
