-- +goose Up
-- +goose StatementBegin
CREATE TABLE products (
    article SERIAL primary key,
    catalog text,
    name text,
    description text, --varchar(255)
    photo_url text NOT NULL default '',
    price FLOAT
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists products;
-- +goose StatementEnd
