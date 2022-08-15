-- +goose Up
-- +goose StatementBegin

CREATE OR REPLACE FUNCTION trigger_set_timestamp()
    RETURNS TRIGGER AS
$$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TABLE files
(
    id         BIGSERIAL               NOT NULL PRIMARY KEY,
    path       TEXT                    NOT NULL UNIQUE,
    name       TEXT                    NOT NULL UNIQUE,

    created_at timestamptz             NOT NULL DEFAULT NOW()
);

CREATE TABLE users
(
    id          BIGSERIAL    NOT NULL PRIMARY KEY,
    email       VARCHAR(100) NOT NULL UNIQUE,
    password    VARCHAR(60),
    username    VARCHAR(60),
    name        VARCHAR(60),
    surname     VARCHAR(60),
    patronymic  VARCHAR(60),

    is_active   BOOLEAN               DEFAULT FALSE,

    token_hash  VARCHAR(15),
    is_verified  BOOLEAN               DEFAULT FALSE,
    is_oauth    BOOLEAN               DEFAULT FALSE,

    avatar_id   BIGINT REFERENCES files,

    created_at  timestamptz  NOT NULL DEFAULT NOW(),
    updated_at  timestamptz  NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX email_idx ON users (email);

CREATE TABLE smers
(
    id         BIGSERIAL               NOT NULL PRIMARY KEY,
    user_id    BIGINT REFERENCES users NOT NULL,

    situation  TEXT                    NOT NULL,
    thoughts   TEXT[]                  NOT NULL,
    emotions   TEXT[]                  NOT NULL,
    reactions  TEXT[]                  NOT NULL,

    created_at timestamptz             NOT NULL DEFAULT NOW(),
    updated_at timestamptz             NOT NULL DEFAULT NOW()
);

CREATE TABLE refresh_tokens
(
    id         BIGSERIAL               NOT NULL PRIMARY KEY,
    user_id    BIGINT REFERENCES users NOT NULL UNIQUE,
    token      TEXT                    NOT NULL,

    created_at timestamptz             NOT NULL DEFAULT NOW()
);

CREATE TRIGGER set_users_timestamp
    BEFORE UPDATE
    ON users
    FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();

CREATE TRIGGER set_smers_timestamp
    BEFORE UPDATE
    ON smers
    FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX email_idx;
DROP TABLE refresh_tokens;
DROP TABLE smers;
DROP TABLE users;
DROP TABLE files;
-- +goose StatementEnd
