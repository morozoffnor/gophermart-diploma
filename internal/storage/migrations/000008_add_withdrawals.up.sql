create table if not exists withdrawals
(
    order_id     varchar(255) not null,
    sum          float        not null,
    user_id      varchar(255) not null,
    withdrawn_at timestamptz default now() not null ,
    constraint withdrawals_orders_id_user_id_fk
        foreign key (order_id) references orders (id)
);

