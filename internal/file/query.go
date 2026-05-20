package file

func getFileQuery() string {
	return `SELECT
			id,
			code,
			token,
			name,
			storage_key,
			visibility,
			is_encrypted,
			size,
			created_at,
			expires_at
		FROM files
		WHERE code = ?
		  AND (expires_at IS NULL OR expires_at > ?)
		LIMIT 1`
}

func createFileQuery() string {
	return `INSERT INTO files (
			id,
			code,
			token,
			name,
			storage_key,
			visibility,
			is_encrypted,
			size,
			created_at,
			expires_at
		)
		VALUES (
			:id,
			:code,
			:token,
			:name,
			:storage_key,
			:visibility,
			:is_encrypted,
			:size,
			:created_at,
			:expires_at
		)`
}

func deleteFileQuery() string {
	return `DELETE FROM files
		WHERE code = ?`
}

func deleteExpiredQuery() string {
	return `DELETE FROM files
	WHERE expires_at IS NOT NULL
	AND expires_at <= ?
	LIMIT ?`
}

func getExpiredStorageKeys() string {
	return `SELECT storage_key FROM files
	WHERE expires_at IS NOT NULL
	AND expires_at <= ?
	LIMIT ?`
}
