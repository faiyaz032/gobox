-- +goose Up
-- +goose StatementBegin
CREATE TABLE box (
    fingerprint_id UUID PRIMARY KEY,
    container_id   TEXT NOT NULL,
    status         TEXT NOT NULL DEFAULT 'starting',
    last_active    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Index for the 24h cleanup worker
CREATE INDEX idx_box_last_active ON box(last_active);
-- Index for status filtering (e.g., checking for 'starting' boxes)
CREATE INDEX idx_box_status ON box(status);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS box;
-- +goose StatementEnd
