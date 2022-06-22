INSERT INTO products_history (id, task_id, article, brand, status, errorresponse, description)
SELECT id, task_id, article, brand, status, errorresponse, description FROM products_buffer
WHERE products_buffer.id NOT IN  (SELECT  id from products_history) AND products_buffer.status = 2;

DELETE FROM products_buffer WHERE status = 2;