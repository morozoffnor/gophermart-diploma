alter table withdrawals
    add constraint withdrawals_orders_id_user_id_fk
        foreign key (order_id) references orders;