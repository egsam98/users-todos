-- +goose Up
create table users (
    id serial not null,
    username varchar not null unique,
    password varchar not null
);

-- +goose Down
drop table users;
