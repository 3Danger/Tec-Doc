INSERT INTO
    tasks (supplier_id, user_id, description)
VALUES
    ((SELECT id FROM suppliers where description = 'IP_IVANOV') ,
     (SELECT id FROM users where supplier_id = (SELECT id FROM suppliers where description = 'IP_IVANOV')),
     'UPLOAD NEW CATALOGUE FOR CANDYSHOP_COMPANY');

INSERT INTO
    tasks (supplier_id, user_id, description)
VALUES
    ((SELECT id FROM suppliers where description = 'IP_PETROV') ,
     (SELECT id FROM users where supplier_id = (SELECT id FROM suppliers where description = 'IP_PETROV')),
     'UPLOAD NEW CATALOGUE FOR AUTOPARTS_COMPANY');

INSERT INTO
    tasks (supplier_id, user_id, description)
VALUES
    ((SELECT id FROM suppliers where description = 'IP_SIDOROV') ,
     (SELECT id FROM users where supplier_id = (SELECT id FROM suppliers where description = 'IP_SIDOROV')),
     'UPLOAD NEW CATALOGUE FOR ROGAIKOPITA_COMPANY');
