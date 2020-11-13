-- +goose Up
create table todos (
    id serial not null,
    title varchar not null,
    description varchar,
    deadline date,
    user_id integer not null
);

-- +goose Down
drop table todos;
