-- +goose Up
create table users (
    id bigint primary key generated always as identity,
    first_name text not null,
    last_name text
);

CREATE TABLE events (
    id bigint primary key generated always as identity,
    name text not null ,
    date_time timestamp not null ,
    duration interval not null check ( duration > interval '0') ,
    description text,
    user_id bigint references users(id) on DELETE cascade ,
    time_in_advance interval CHECK (time_in_advance >= INTERVAL '0')
);

create index on users(first_name, last_name);

create index on events(date_time);

create index on events(user_id);

-- +goose Down
drop table users;
drop table events;