-- +goose Up
-- +goose StatementBegin
CREATE TABLE session_containers (
    id SERIAL PRIMARY KEY,
    session_id   VARCHAR(256) NOT NULL,
    container_id VARCHAR(256) NOT NULL,
    paused_at    TIMESTAMPTZ NULL,
    is_paused    BOOLEAN DEFAULT false
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE session_containers;
-- +goose StatementEnd
