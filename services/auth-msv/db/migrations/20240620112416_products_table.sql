-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS Users (
    id UUID NOT NULL,
    login VARCHAR(128) NOT NULL,
    password VARCHAR(64) NOT NULL,
    access_level INTEGER NOT NULL,
    PRIMARY KEY (id),
    UNIQUE (login)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS Users;
-- +goose StatementEnd
