-- +goose Up
-- +goose StatementBegin
DO
$$
    BEGIN
        CREATE TYPE scope_type AS ENUM ('authentication', 'verification');
    EXCEPTION
        WHEN duplicate_object THEN null;
    END
$$;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP type IF EXISTS scope_type;
-- +goose StatementEnd
