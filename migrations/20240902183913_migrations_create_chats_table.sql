-- +goose Up
-- +goose StatementBegin
CREATE TYPE text_array AS ARRAY OF text;

CREATE TABLE
    IF NOT EXISTS chats (
        id uuid PRIMARY KEY DEFAULT gen_random_uuid (),
        usernames text_array
    );

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TYPE IF EXISTS text_array;

DROP TABLE IF EXISTS chats;

-- +goose StatementEnd