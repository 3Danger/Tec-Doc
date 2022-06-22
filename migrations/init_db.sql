CREATE TABLE IF NOT EXISTS suppliers (
     id             BIGSERIAL PRIMARY KEY,
--      user_id        BIGINT,
     upload_date    TIMESTAMPTZ,
     update_date    TIMESTAMPTZ,
     status         INT,
     description    TEXT
);

CREATE TABLE IF NOT EXISTS users (
     id             BIGSERIAL PRIMARY KEY,
     supplier_id    BIGINT REFERENCES suppliers ON DELETE CASCADE ,
     description    TEXT
);

CREATE TABLE IF NOT EXISTS tasks (
     id             BIGSERIAL PRIMARY KEY,
     supplier_id    BIGINT REFERENCES suppliers ON DELETE CASCADE ,
     user_id        BIGINT REFERENCES users ON DELETE CASCADE ,
     description    text
);

CREATE TABLE IF NOT EXISTS products_buffer (
   id             BIGSERIAL PRIMARY KEY,
   upload_id      BIGINT REFERENCES tasks ON DELETE CASCADE ,
   article        TEXT,
   brand          TEXT,
   status         INT,
   errorResponse  TEXT,
   description    TEXT
);

CREATE TABLE IF NOT EXISTS products_history (
    id             BIGSERIAL PRIMARY KEY,
    upload_id      BIGINT REFERENCES tasks ON DELETE CASCADE,
    article        TEXT,
    brand          TEXT,
    status         INT,
    errorResponse  TEXT,
    description    TEXT
);
