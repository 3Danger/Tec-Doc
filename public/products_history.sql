create table products_history
(
    id            bigserial
        primary key,
    task_id       bigint,
    article       text,
    brand         text,
    status        integer,
    errorresponse text,
    description   varchar
);

alter table products_history
    owner to tecdoc;

