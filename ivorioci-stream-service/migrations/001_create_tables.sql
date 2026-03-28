-- Stream service tables
-- Run this once against the shared PostgreSQL database (ivorioci).
-- The gateway manages its own tables (users, refresh_tokens) via Prisma.

CREATE TABLE IF NOT EXISTS categories (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name        VARCHAR(255) NOT NULL,
    slug        VARCHAR(255) NOT NULL UNIQUE,
    description TEXT         NOT NULL DEFAULT '',
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS videos (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title         VARCHAR(255) NOT NULL,
    description   TEXT         NOT NULL DEFAULT '',
    -- file_path is relative to VIDEO_STORAGE_PATH (e.g. "documentaries/ocean.mp4")
    file_path     VARCHAR(500) NOT NULL,
    thumbnail_url VARCHAR(500) NOT NULL DEFAULT '',
    category_id   UUID         REFERENCES categories(id) ON DELETE SET NULL,
    duration      INTEGER      NOT NULL DEFAULT 0,  -- seconds
    views_count   INTEGER      NOT NULL DEFAULT 0,
    file_size     BIGINT       NOT NULL DEFAULT 0,  -- bytes
    mime_type     VARCHAR(100) NOT NULL DEFAULT 'video/mp4',
    is_published  BOOLEAN      NOT NULL DEFAULT true,
    created_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_videos_category_id  ON videos(category_id);
CREATE INDEX IF NOT EXISTS idx_videos_is_published ON videos(is_published);
CREATE INDEX IF NOT EXISTS idx_videos_created_at   ON videos(created_at DESC);

-- Auto-update updated_at on row change
CREATE OR REPLACE FUNCTION update_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_categories_updated_at ON categories;
CREATE TRIGGER trg_categories_updated_at
    BEFORE UPDATE ON categories
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();

DROP TRIGGER IF EXISTS trg_videos_updated_at ON videos;
CREATE TRIGGER trg_videos_updated_at
    BEFORE UPDATE ON videos
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();
