alter table orders
    alter column accrual type integer using accrual::integer;
