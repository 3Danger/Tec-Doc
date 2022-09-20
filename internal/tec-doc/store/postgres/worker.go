package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v4"
	"tec-doc/pkg/model"
	"time"
)

func (s *store) GetOldestTask(ctx context.Context, tx Transaction) (int64, string, error) {
	var (
		getOldestTaskQuery = `SELECT id, supplier_id_string FROM tasks.tasks  WHERE status=$1 or status=$2 ORDER BY upload_date ASC LIMIT 1;`
		executor           Executor
		t                  model.Task
	)

	executor = s.pool
	if tx != nil {
		executor = tx
	}

	ctx, cancel := context.WithTimeout(ctx, s.cfg.Timeout)
	defer cancel()

	row := executor.QueryRow(ctx, getOldestTaskQuery, StatusNew, StatusProcess)
	if err := row.Scan(&t.ID, &t.SupplierIdString); err != nil {
		return 0, "", err
	}
	return t.ID, t.SupplierIdString, nil
}

func (s *store) GetProductsBufferWithStatus(ctx context.Context, tx Transaction, uploadID int64, limit int, offset int, status int) ([]model.Product, error) {
	var (
		executor       Executor
		productsBuffer = make([]model.Product, 0)
	)

	query := `
	SELECT 
			id, 
			upload_id, 
			article, 
			article_supplier, 
			brand,
			barcode,
			subject,
	COALESCE(price, 0),
			upload_date,
			update_date,
			status,
	COALESCE(errorresponse, '')
	FROM tasks.products_buffer WHERE upload_id = $1 and status=$2 ORDER BY update_date LIMIT $3 OFFSET $4`

	executor = s.pool

	if tx != nil {
		executor = tx
	}

	ctx, cancel := context.WithTimeout(ctx, s.cfg.Timeout)
	defer cancel()

	rows, err := executor.Query(ctx, query, uploadID, status, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("can't get products from buffer: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var p model.Product
		err = rows.Scan(&p.ID, &p.UploadID, &p.Article, &p.ArticleSupplier, &p.Brand, &p.Barcode, &p.Subject, &p.Price, &p.UploadDate, &p.UpdateDate, &p.Status, &p.ErrorResponse)
		if err != nil {
			return nil, fmt.Errorf("can't get products from buffer: %w", err)
		}
		productsBuffer = append(productsBuffer, p)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("can't get products from buffer: %w", err)
	}

	return productsBuffer, nil
}

func (s *store) UpdateProductStatus(ctx context.Context, tx Transaction, productID int64, status int) error {

	var (
		query             = `UPDATE tasks.products_buffer SET status=$1, update_date = $2 WHERE id=$3;`
		executor Executor = s.pool
	)

	if tx != nil {
		executor = tx
	}

	ctx, cancel := context.WithTimeout(ctx, s.cfg.Timeout)
	defer cancel()

	res, err := executor.Exec(ctx, query, status, time.Now().UTC(), productID)
	if err != nil {
		return fmt.Errorf("can't exec updateProductStatusQuery: %w", err)
	}

	if res.RowsAffected() == 0 {
		return fmt.Errorf("no rows were updated")
	}

	return nil
}

func (s *store) UpdateTaskProductsNumber(ctx context.Context, tx Transaction, uploadID, productsFailed, productsProcessed int64) error {
	var (
		updateTaskProductsNumberQuery = `UPDATE tasks.tasks SET products_failed=products_failed + $1, products_processed=products_processed + $2, update_date = $3
                  WHERE id=$4;`
		executor Executor
	)

	executor = s.pool
	if tx != nil {
		executor = tx
	}

	ctx, cancel := context.WithTimeout(ctx, s.cfg.Timeout)
	defer cancel()

	res, err := executor.Exec(ctx, updateTaskProductsNumberQuery, productsFailed, productsProcessed, time.Now().UTC(), uploadID)
	if err != nil {
		return fmt.Errorf("can't exec updateTaskProductsNumberQuery: %w", err)
	}

	if res.RowsAffected() == 0 {
		return fmt.Errorf("no rows were updated")
	}
	return nil
}

func (s *store) UpdateTaskStatus(ctx context.Context, tx Transaction, uploadID int64, status int) error {
	var (
		updateTaskStatusQuery          = `UPDATE tasks.tasks SET status=$1, update_date = $2 WHERE id=$3;`
		executor              Executor = s.pool
	)

	if tx != nil {
		executor = tx
	}

	ctx, cancel := context.WithTimeout(ctx, s.cfg.Timeout)
	defer cancel()

	res, err := executor.Exec(ctx, updateTaskStatusQuery, status, time.Now().UTC(), uploadID)
	if err != nil {
		return fmt.Errorf("can't exec updateTaskStatusQuery: %w", err)
	}

	if res.RowsAffected() == 0 {
		return fmt.Errorf("no rows were updated")
	}

	return nil
}

func (s *store) MoveProductsToHistoryByUploadId(ctx context.Context, tx Transaction, uploadId int64) (err error) {
	var executor Executor = s.pool
	if tx != nil {
		executor = tx
	}
	if _, err = executor.Exec(ctx, "CALL tasks.move_products_from_buffer_to_history($1)", uploadId); err != nil {
		return err
	}
	return nil
}

func (s *store) UpdateProductBuffer(ctx context.Context, tx Transaction, products *model.Product) (err error) {
	var executor Executor = s.pool
	if tx != nil {
		executor = tx
	}
	ctx, cancel := context.WithTimeout(ctx, s.cfg.Timeout)
	defer cancel()

	query :=
		`UPDATE tasks.products_buffer SET status = $1, errorresponse = $2, update_date = $3 WHERE id = $4`
	row := executor.QueryRow(ctx, query, products.Status, products.ErrorResponse, time.Now().UTC(), products.ID)
	if err = row.Scan(); err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return err
	}
	return nil
}
