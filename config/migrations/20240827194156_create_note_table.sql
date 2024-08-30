-- +goose Up
CREATE TABLE note (
    id serial primary key,
    title text not null,
    body text not null,
    created_at timestamp not null default now(),
    updated_at timestamp
);

-- +goose Down
DROP TABLE note;
