-- +goose Up
-- +goose StatementBegin
CREATE TABLE orders (
    id SERIAL primary key,
    id_user integer,
    status_order integer,
    pvz jsonb,
    orderr text,
    CREATED_AT timestamp NOT NULL DEFAULT NOW() at time zone 'utc'
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists orders;
-- +goose StatementEnd
