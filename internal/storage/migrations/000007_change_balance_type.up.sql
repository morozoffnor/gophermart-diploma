alter table users
    alter column balance type float using balance::float;
alter table users
    add withdrawn float default 0 not null;