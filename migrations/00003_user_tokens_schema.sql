-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS user_tokens (
    token_hash bytea NOT NULL,
    user_id uuid NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    scope scope_type NOT NULL,
    PRIMARY KEY (token_hash),
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS user_tokens;

-- +goose StatementEnd
