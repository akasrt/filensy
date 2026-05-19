-- +goose Up
CREATE TABLE IF NOT EXISTS files (
    id           VARCHAR(36) PRIMARY KEY,
    code         VARCHAR(8) NOT NULL,
    token        CHAR(64) NOT NULL,
    name         VARCHAR(255) NOT NULL,
    storage_key  VARCHAR(255) NOT NULL,
    visibility    ENUM('private', 'public') NOT NULL DEFAULT 'private',
    is_encrypted  BOOLEAN NOT NULL DEFAULT FALSE,
    size         BIGINT UNSIGNED NOT NULL,
    created_at   TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expires_at   TIMESTAMP NULL DEFAULT NULL,

    UNIQUE KEY uq_file_code (code),
    UNIQUE KEY uq_file_storage_Key(storage_key),

    INDEX idx_expiry (expires_at)
);

-- +goose Down

