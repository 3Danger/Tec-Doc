create table products_buffer
(
    id            bigserial
        primary key,
    task_id       bigint
        references tasks
            on delete cascade,
    article       text,
    brand         text,
    status        integer,
    errorresponse text,
    description   text
);

alter table products_buffer
    owner to tecdoc;

