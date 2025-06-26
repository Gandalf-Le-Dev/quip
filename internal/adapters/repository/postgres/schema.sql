CREATE TABLE IF NOT EXISTS files (
    id VARCHAR(11) PRIMARY KEY,
    original_name VARCHAR(255) NOT NULL,
    size BIGINT NOT NULL,
    content_type VARCHAR(100) NOT NULL,
    storage_key VARCHAR(100) NOT NULL UNIQUE,
    downloads INT NOT NULL DEFAULT 0,
    max_downloads INT NOT NULL DEFAULT -1,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS pastes (
    id VARCHAR(11) PRIMARY KEY,
    content TEXT NOT NULL,
    language VARCHAR(50) NOT NULL,
    title VARCHAR(255),
    views INT NOT NULL DEFAULT 0,
    max_views INT NOT NULL DEFAULT -1,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMP NOT NULL
);

CREATE INDEX idx_files_expires_at ON files(expires_at);
CREATE INDEX idx_pastes_expires_at ON pastes(expires_at);