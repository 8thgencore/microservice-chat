-- +goose Up
-- +goose StatementBegin
CREATE TABLE
    IF NOT EXISTS transaction_log (
        id uuid primary key DEFAULT gen_random_uuid(),
        timestamp timestamp not null default now (),
        log text not null
    );

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS transaction_log;

-- +goose StatementEnd