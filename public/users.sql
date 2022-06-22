create table users
(
    id          bigserial
        primary key,
    supplier_id bigint
        references suppliers
            on delete cascade,
    description text
);

alter table users
    owner to tecdoc;

