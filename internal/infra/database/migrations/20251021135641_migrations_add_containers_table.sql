-- +goose Up
-- +goose StatementBegin
CREATE TABLE session_containers (
    id SERIAL PRIMARY KEY,
    session_id   VARCHAR(256) NOT NULL,
    container_id VARCHAR(256) NOT NULL,
    last_active  TIMESTAMP NOT NULL
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE session_containers
-- +goose StatementEnd
