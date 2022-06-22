INSERT INTO products_buffer
    (task_id, article, brand, status, errorresponse, description)
VALUES
    ((SELECT id from tasks where description = 'UPLOAD NEW CATALOGUE FOR CANDYSHOP_COMPANY'),
     '010001N',
     'AKS DASIS',
     0,
     'invalid product_info',
     'Радиатор, охлаждение двигателя');

INSERT INTO products_buffer
(task_id, article, brand, status, errorresponse, description)
VALUES
    ((SELECT id from tasks where description = 'UPLOAD NEW CATALOGUE FOR CANDYSHOP_COMPANY'),
     '0092L40270',
     'BOSCH',
     0,
     'invalid product_info',
     'Стартерная аккумуляторная батарея');

INSERT INTO products_buffer
(task_id, article, brand, status, errorresponse, description)
VALUES
    ((SELECT id from tasks where description = 'UPLOAD NEW CATALOGUE FOR CANDYSHOP_COMPANY'),
     '000.001-00A',
     'PE Automotive',
     0,
     'invalid product_info',
     'Брызговик');

INSERT INTO products_buffer
(task_id, article, brand, status, errorresponse, description)
VALUES
    ((SELECT id from tasks where description = 'UPLOAD NEW CATALOGUE FOR CANDYSHOP_COMPANY'),
     '10-3011',
     'Airstal',
     0,
     'invalid product_info',
     'Компрессор, кондиционер');

INSERT INTO products_buffer
(task_id, article, brand, status, errorresponse, description)
VALUES
    ((SELECT id from tasks where description = 'UPLOAD NEW CATALOGUE FOR CANDYSHOP_COMPANY'),
     '103011',
     'AUTOGAMMA',
     0,
     'invalid product_info',
     'Интеркулер');

INSERT INTO products_buffer
(task_id, article, brand, status, errorresponse, description)
VALUES
    ((SELECT id from tasks where description = 'UPLOAD NEW CATALOGUE FOR AUTOPARTS_COMPANY'),
     '103002',
     'FARCOM',
     0,
     'invalid product_info',
     'Стартер');

INSERT INTO products_buffer
(task_id, article, brand, status, errorresponse, description)
VALUES
    ((SELECT id from tasks where description = 'UPLOAD NEW CATALOGUE FOR AUTOPARTS_COMPANY'),
     '0451103274',
     'BOSCH',
     0,
     'invalid product_info',
     'Масляный фильтр');
