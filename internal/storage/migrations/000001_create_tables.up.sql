create table if not exists users
(
    id       varchar(255) not null
        constraint users_pk
            primary key,
    login    varchar(255) not null,
    password varchar(255) not null
);

create table if not exists orders
(
    id      integer not null
        constraint orders_pk
            primary key,
    a       integer,
    b       integer,
    user_id varchar(255)
        constraint orders_users_id_fk
            references users
);

