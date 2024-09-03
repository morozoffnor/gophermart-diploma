alter table orders
    alter column id type varchar(255) using id::varchar(255);