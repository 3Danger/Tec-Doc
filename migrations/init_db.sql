CREATE SCHEMA IF NOT EXISTS tasks;

CREATE TABLE IF NOT EXISTS tasks.tasks (
     id                 BIGSERIAL PRIMARY KEY,
     supplier_id        BIGINT,
     user_id            BIGINT,
     IP                 TEXT,
     upload_date        TIMESTAMPTZ,
     update_date        TIMESTAMPTZ,
     status             INT,
     products_processed INT,
     products_failed    INT,
     products_total     INT
);

CREATE INDEX ON tasks.tasks(supplier_id);

CREATE TABLE IF NOT EXISTS tasks.products_buffer (
   id                   BIGSERIAL PRIMARY KEY,
   upload_id            BIGINT,
   article              TEXT,
   article_supplier     TEXT,
   brand                TEXT,
   barcode              TEXT,
   subject              TEXT,
   price                INT,
   upload_date          TIMESTAMPTZ,
   update_date          TIMESTAMPTZ,
   amount               INT,
   status               INT,
   errorResponse        TEXT
);

CREATE INDEX ON tasks.products_buffer(upload_id);

CREATE TABLE IF NOT EXISTS tasks.products_history (
   id                   BIGSERIAL PRIMARY KEY,
   upload_id            BIGINT,
   article              TEXT,
   article_supplier     TEXT,
   brand                TEXT,
   barcode              TEXT,
   subject              TEXT,
   price                INT,
   upload_date          TIMESTAMPTZ,
   update_date          TIMESTAMPTZ,
   amount               INT,
   status               INT,
   errorResponse        TEXT
);

CREATE INDEX ON tasks.products_history(upload_id);
