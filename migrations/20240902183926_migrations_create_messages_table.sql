-- +goose Up
-- +goose StatementBegin
CREATE TABLE
    IF NOT EXISTS messages (
        chat_id uuid references chats (id),
        from_user text,
        text text,
        timestamp timestamp not null
    );

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS messages;

-- +goose StatementEnd