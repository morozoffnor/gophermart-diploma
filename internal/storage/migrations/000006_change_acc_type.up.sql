alter table orders
    alter column accrual type float using accrual::float;