create table suppliers
(
    id          bigserial
        primary key,
    user_id     bigint,
    upload_date timestamp with time zone,
    update_date timestamp with time zone,
    status      integer,
    description text
);

alter table suppliers
    owner to tecdoc;

