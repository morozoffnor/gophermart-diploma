alter table users
    alter column balance type integer using balance::integer;
alter table users
    drop column withdrawn;
