-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS projects(
    id uuid NOT NULL,
    name TEXT NOT NULL,
    description TEXT,
    workspace_id uuid NOT NULL,
    start_date DATE,
    end_date DATE,
    created_at TIMESTAMP DEFAULT now() NOT NULL,
    last_modified TIMESTAMP DEFAULT now() NOT NULL,
    PRIMARY KEY (id),
    FOREIGN KEY (workspace_id) REFERENCES workspaces(id)    
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS projects;
-- +goose StatementEnd
