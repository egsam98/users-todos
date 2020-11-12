create table if not exists users (
    id serial,
    username varchar not null unique,
    password varchar not null
);

-- name: CreateUser :exec
insert into users (username, password) values ($1, $2);