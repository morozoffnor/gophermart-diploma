alter table orders
    alter column status type varchar(20) using status::varchar(20);