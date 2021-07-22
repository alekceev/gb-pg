-- список листов TODO

create database todo;

create user www password 'pass';

grant all privileges on database todo to www;

CREATE extension IF NOT EXISTS "uuid-ossp";

create table users (
    uuid uuid not null default uuid_generate_v4() primary key,
    ctime timestamp not null default current_timestamp,
    name varchar(127),
    email varchar(127) unique,
    pass varchar(32),
    salt varchar(10)
);

insert into users (name, email, pass, salt) values
  ('User1', 'user1@mail.ru', 'd36f9d30acb4e2857a2818aa8420f7b7', '111'),
  ('User2', 'user2@mail.ru', 'd36f9d30acb4e2857a2818aa8420f7b7', '111'),
  ('Admin', 'admin@mail.ru', '66e1a360ee8070ba822aca90526dec47', '222');
