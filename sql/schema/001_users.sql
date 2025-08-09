-- +goose Up
CREATE TABLE users (
    id UUID PRIMARY KEY, -- https://www.postgresql.org/docs/16/datatype-uuid.html
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    email TEXT NOT NULL,
    UNIQUE(  email )
);  -- https://www.postgresql.org/docs/16/sql-createtable.html

-- +goose Down
DROP TABLE users;