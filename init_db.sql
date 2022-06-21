CREATE TABLE IF NOT EXISTS suppliers (
     id             bigserial primary key,
     user_id        bigint,
     upload_date    timestamptz,
     update_date    timestamptz,
     status         int,
     description    text
);
CREATE TABLE IF NOT EXISTS users (
     id             bigserial primary key,
     supplier_id    bigint references suppliers on delete cascade,
     description    text
);
CREATE TABLE IF NOT EXISTS tasks (
     id             bigserial primary key,
     supplier_id    bigint references suppliers on delete cascade ,
     user_id        bigint references users on delete cascade ,
     description    text
);
CREATE TABLE IF NOT EXISTS products_buffer (
     id             bigserial primary key,
     task_id        bigint references tasks on delete cascade ,
     article        text,
     brand          text,
     status         int,
     errorResponse  text,
     description    text
);
CREATE TABLE IF NOT EXISTS products_history (
    id             bigserial primary key,
    task_id        bigint references tasks on delete cascade ,
    article        text,
    brand          text,
    status         int,
    errorResponse  text,
    description    text
);
