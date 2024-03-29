-- список листов TODO

-- таблица юзеров
create table if not exists users (
    uuid uuid not null default uuid_generate_v4() primary key,
    ctime     timestamp not null default current_timestamp,
    name      varchar(127),
    email     varchar(127) unique,
    pass      varchar(32),
    salt      varchar(10)
);

-- список туду листов юзера
create table if not exists lists (
    id          serial not null unique,
    user_uuid   uuid references users (uuid) on delete cascade not null,
    title       varchar(255) not null,
    description varchar(255)
);

-- список задач
create table if not exists items (
    id          serial       not null unique,
    list_id     int references lists (id) on delete cascade not null,
    title       varchar(255) not null,
    description varchar(255),
    due_date    timestamp not null,
    done        boolean      not null default false
);
