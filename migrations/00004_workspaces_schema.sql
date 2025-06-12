-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS workspaces(
    id uuid PRIMARY KEY,
    name VARCHAR(120) NOT NULL,
    description TEXT,
    user_id uuid NOT NULL,
    created_at TIMESTAMP DEFAULT now() NOT NULL,
    last_modified TIMESTAMP DEFAULT now() NOT NULL,
    FOREIGN KEY(user_id) REFERENCES users(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS workspaces;
-- +goose StatementEnd
