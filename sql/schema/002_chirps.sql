-- +goose Up
CREATE TABLE chirps (
    id UUID PRIMARY KEY, -- https://www.postgresql.org/docs/16/datatype-uuid.html
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    body TEXT NOT NULL,
    user_id UUID NOT NULL
        REFERENCES users (id) ON DELETE CASCADE
);  -- https://www.postgresql.org/docs/16/sql-createtable.html

-- +goose Down
DROP TABLE chirps;