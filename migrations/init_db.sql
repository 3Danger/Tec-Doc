CREATE SCHEMA IF NOT EXISTS tasks;

CREATE TABLE IF NOT EXISTS tasks.tasks (
     id                 BIGSERIAL PRIMARY KEY,
     supplier_id_string VARCHAR(64) not null,
     supplier_id        BIGINT not null,
     user_id            BIGINT not null,
     IP                 TEXT not null,
     upload_date        TIMESTAMPTZ not null,
     update_date        TIMESTAMPTZ not null,
     status             INT not null,
     products_processed INT not null,
     products_failed    INT not null,
     products_total     INT not null
);

CREATE INDEX ON tasks.tasks(supplier_id);

CREATE TABLE IF NOT EXISTS tasks.products_buffer (
   id                   BIGSERIAL PRIMARY KEY not null,
   upload_id            BIGINT not null,
   article              TEXT not null,
   article_supplier     TEXT not null,
   brand                TEXT not null,
   barcode              TEXT not null,
   subject              TEXT not null,
   price                INT,
   upload_date          TIMESTAMPTZ not null,
   update_date          TIMESTAMPTZ not null,
   amount               INT not null,
   status               INT not null,
   errorResponse        TEXT
);

CREATE INDEX ON tasks.products_buffer(upload_id);

CREATE TABLE IF NOT EXISTS tasks.products_history (
   id                   BIGSERIAL PRIMARY KEY not null,
   upload_id            BIGINT not null,
   article              TEXT not null,
   article_supplier     TEXT not null,
   brand                TEXT not null,
   barcode              TEXT not null,
   subject              TEXT not null,
   price                INT,
   upload_date          TIMESTAMPTZ not null,
   update_date          TIMESTAMPTZ not null,
   amount               INT not null,
   status               INT not null,
   errorResponse        TEXT
);

CREATE INDEX ON tasks.products_history(upload_id);
