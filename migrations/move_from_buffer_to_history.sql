INSERT INTO products_history (id, upload_id, article, brand, status, errorResponse, description)
    SELECT id, upload_id, article, brand, status, errorResponse, description FROM products_buffer
WHERE products_buffer.id NOT IN  (SELECT  id from products_history)
  AND products_buffer.status = 0 AND products_buffer.upload_id = 1;

DELETE FROM products_buffer WHERE status = 2 AND products_buffer.upload_id = 1;