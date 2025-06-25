-- name: CreateFile :one
INSERT INTO files (
    id, original_name, size, content_type, storage_key,
    downloads, max_downloads, created_at, expires_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9
) RETURNING *;

-- name: GetFileByID :one
SELECT * FROM files WHERE id = $1 LIMIT 1;

-- name: IncrementFileDownloads :exec
UPDATE files SET downloads = downloads + 1 WHERE id = $1;

-- name: DeleteExpiredFiles :exec
DELETE FROM files WHERE expires_at < NOW();

-- name: CreatePaste :one
INSERT INTO pastes (
    id, content, language, title, views, max_views, created_at, expires_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
) RETURNING *;

-- name: GetPasteByID :one
SELECT * FROM pastes WHERE id = $1 LIMIT 1;

-- name: IncrementPasteViews :exec
UPDATE pastes SET views = views + 1 WHERE id = $1;

-- name: DeleteExpiredPastes :exec
DELETE FROM pastes WHERE expires_at < NOW();