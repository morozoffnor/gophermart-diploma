alter table orders
         rename column a to status;
alter table orders
         alter column status set not null;
alter table orders
         rename column b to accrual;
alter table orders
         alter column user_id set not null;
set timezone to 'Europe/Moscow';
alter table orders
         add uploaded_at timestamptz not null default now()