-- +goose Up
create table todos (
    id serial not null,
    title varchar not null,
    description varchar,
    deadline date
);

-- +goose Down
drop table todos;
