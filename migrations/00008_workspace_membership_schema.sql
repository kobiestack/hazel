-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS workspace_memberships(
    workspace_id uuid NOT NULL,
    user_id uuid NOT NULL,
    FOREIGN KEY (workspace_id) REFERENCES workspaces(id),
    FOREIGN KEY (user_id) REFERENCES users(id),
    PRIMARY KEY (workspace_id, user_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS workspace_memberships;
-- +goose StatementEnd
