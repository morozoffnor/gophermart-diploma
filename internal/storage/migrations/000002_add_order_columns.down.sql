alter table orders
         rename column status to a;
alter table orders
         alter column status drop not null ;
alter table orders
         rename column accrual to b;
alter table orders
         alter column user_id drop not null;
set time zone 'Europe/Moscow';
alter table orders
    drop column uploaded_at;