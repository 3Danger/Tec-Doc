create table tasks
(
    id          bigserial
        primary key,
    supplier_id bigint
        references suppliers
            on delete cascade,
    user_id     bigint
        references users
            on delete cascade,
    description text
);

alter table tasks
    owner to tecdoc;

