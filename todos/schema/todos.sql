create table if not exists todos (
    id serial not null,
    title varchar not null,
    description varchar,
    deadline date,
    user_id integer not null
);

-- name: CreateTodo :exec
insert into todos (title, description, deadline, user_id) values ($1, $2, $3, $4);