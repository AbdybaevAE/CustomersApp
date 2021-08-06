create table if not exists customers(
    customer_id serial not null primary key,
    customer_first_name varchar(150) not null,
    customer_last_name varchar(150) not null,
    customer_gender varchar(10),
    customer_email varchar(150) not null unique,
    customer_address varchar(300),
    customer_birth_date date not null,
    customer_created_at timestamp not null default now(),
    customer_updated_at timestamp not null default now(),
    customer_hash varchar(20) not null
)