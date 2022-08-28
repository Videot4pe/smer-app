-- +goose Up
-- +goose StatementBegin
ALTER TABLE users
ADD hash_expire timestamptz DEFAULT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users
DROP COLUMN hash_expire;
-- +goose StatementEnd
