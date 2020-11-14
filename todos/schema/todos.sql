create table if not exists todos (
    id serial not null,
    title varchar not null,
    description varchar,
    deadline timestamp,
    user_id integer not null
);

-- name: CreateTodo :one
insert into todos (title, description, deadline, user_id) values ($1, $2, $3, $4) returning *;

-- name: FindTodoById :one
select * from todos where id = $1;

-- name: UpdateTodo :one
update todos set title=$1, description=$2, deadline=$3 where id=$4 returning *;

-- name: DeleteTodo :execrows
delete from todos where id = $1 and user_id = $2;

-- name: FindAll :many
select * from todos where user_id = $1 order by deadline;

-- name: FindBeforeDeadline :many
select * from todos where deadline < $1 and user_id = $2 order by deadline;