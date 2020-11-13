create table if not exists users (
    id serial not null,
    username varchar not null unique,
    password varchar not null
);

-- name: CreateUser :exec
insert into users (username, password) values ($1, $2);

-- name: FindUser :one
select * from users where username = $1 and password = $2 limit 1;

-- name: FindUserById :one
select * from users where id = $1;