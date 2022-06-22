INSERT INTO
    users (supplier_id, description)
VALUES
    ((SELECT id FROM suppliers WHERE description = 'IP_IVANOV'), 'MANAGER SVETLANA');

INSERT INTO
    users (supplier_id, description)
VALUES
    ((SELECT id FROM suppliers WHERE description = 'IP_PETROV'), 'MANAGER OLEG');

INSERT INTO
    users (supplier_id, description)
VALUES
    ((SELECT id FROM suppliers WHERE description = 'IP_SIDOROV'), 'MANAGER BORIS');