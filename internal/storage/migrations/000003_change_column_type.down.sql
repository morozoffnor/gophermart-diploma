alter table orders
    alter column id type integer using id::integer;