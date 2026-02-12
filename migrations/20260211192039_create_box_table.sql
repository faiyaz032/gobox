-- +goose Up
-- +goose StatementBegin
-- Create enum type for status
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'box_status') THEN
        CREATE TYPE box_status AS ENUM ('running', 'paused');
    END IF;
END$$;

CREATE TABLE box (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    fingerprint_id TEXT NOT NULL,
    container_id   TEXT NOT NULL,
    status         box_status NOT NULL DEFAULT 'running',
    last_active    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Index for the 24h cleanup worker
CREATE INDEX idx_box_last_active ON box(last_active);
-- Index for status filtering
CREATE INDEX idx_box_status ON box(status);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS box;
DROP TYPE IF EXISTS box_status;
-- +goose StatementEnd
