CREATE OR REPLACE PROCEDURE tasks.move_products_from_buffer_to_history(
    _upload_id bigint)
    security definer
    language plpgsql
AS
$$
BEGIN
    WITH tmp AS
        (
            DELETE FROM tasks.products_buffer AS b
                WHERE upload_id = _upload_id
                RETURNING b.id, b.upload_id, b.article, b.article_supplier, b.brand, b.barcode, b.subject, b.price, b.upload_date, b.update_date, b.amount, b.status, b.errorresponse
        )
    INSERT INTO tasks.products_history (id, upload_id, article, article_supplier, brand, barcode, subject, price, upload_date, update_date, amount, status, errorresponse)
        SELECT id, upload_id, article, article_supplier, brand, barcode, subject, price, upload_date, update_date, amount, status, errorresponse
        FROM tmp;
END;
$$;