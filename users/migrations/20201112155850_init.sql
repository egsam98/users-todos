-- +goose Up
create table if not exists users (
    id serial,
    username varchar not null unique,
    password varchar not null
);

-- +goose Down
drop table users;
