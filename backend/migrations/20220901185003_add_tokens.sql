-- +goose Up
-- +goose StatementBegin

-- Refresh tokens will be moved to common tokens table.
DROP TABLE refresh_tokens;

ALTER TABLE users
DROP COLUMN token_hash;

-- New common table for tokens.
CREATE TABLE tokens
(
    user_id    BIGINT REFERENCES users NOT NULL,
    token      TEXT                    NOT NULL,
    token_type VARCHAR(20)             NOT NULL, -- ACTIVATE | AUTH | RESET_PASSWORD
    
    created_at timestamptz             NOT NULL DEFAULT NOW(),

    PRIMARY KEY(user_id, token_type)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE tokens;

ALTER TABLE users
ADD COLUMN token_hash VARCHAR(15);

CREATE TABLE refresh_tokens
(
    id         BIGSERIAL               NOT NULL PRIMARY KEY,
    user_id    BIGINT REFERENCES users NOT NULL UNIQUE,
    token      TEXT                    NOT NULL,

    created_at timestamptz             NOT NULL DEFAULT NOW()
);
-- +goose StatementEnd
