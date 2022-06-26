CREATE TABLE IF NOT EXISTS tasks (
     id                 BIGSERIAL PRIMARY KEY,
     supplier_id        BIGINT,
     user_id            BIGINT,
     upload_date        TIMESTAMPTZ,
     update_date        TIMESTAMPTZ,
     IP                 TEXT,
     status             INT,
     products_processed INT,
     products_failed    INT,
     products_total     INT
);

CREATE TABLE IF NOT EXISTS products_buffer (
   id             BIGSERIAL PRIMARY KEY,
   upload_id      BIGINT,
   article        TEXT,
   brand          TEXT,
   upload_date    TIMESTAMPTZ,
   update_date    TIMESTAMPTZ,
   status         INT,
   errorResponse  TEXT
);

CREATE TABLE IF NOT EXISTS products_history (
    id             BIGSERIAL PRIMARY KEY,
    upload_id      BIGINT,
    article        TEXT,
    brand          TEXT,
    upload_date    TIMESTAMPTZ,
    update_date    TIMESTAMPTZ,
    status         INT,
    errorResponse  TEXT
);
