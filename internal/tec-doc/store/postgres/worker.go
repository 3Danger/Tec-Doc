package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"tec-doc/pkg/model"
	"time"
)

func (s *store) GetOldestTask(ctx context.Context, tx Transaction) (int64, error) {
	var (
		getOldestTaskQuery = `SELECT id FROM tasks.tasks  WHERE status=$1 or status=$2 ORDER BY upload_date ASC LIMIT 1;`
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
	err := row.Scan(&t.ID)
	return t.ID, err
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
	COALESCE(article, ''), 
	COALESCE(article_supplier, ''), 
	COALESCE(brand, ''),
	COALESCE(barcode, ''),
	COALESCE(subject, ''),
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
		updateProductStatusQuery = `UPDATE tasks.products_buffer SET status=$1, update_date = $2 WHERE id=$3;`
		executor                 Executor
	)

	executor = s.pool
	if tx != nil {
		executor = tx
	}

	ctx, cancel := context.WithTimeout(ctx, s.cfg.Timeout)
	defer cancel()

	res, err := executor.Exec(ctx, updateProductStatusQuery, status, time.Now().UTC(), productID)
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
		updateTaskStatusQuery = `UPDATE tasks.tasks SET status=$1, update_date = $2 WHERE id=$3;`
		executor              Executor
	)

	executor = s.pool
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
	var _tx = tx
	if tx == nil {
		if _tx, err = s.pool.Begin(ctx); err != nil {
			return err
		}
		defer func() { _ = _tx.Rollback(ctx) }()
	}
	ctx, cancel := context.WithTimeout(ctx, s.cfg.Timeout)
	defer cancel()

	query := `
INSERT INTO tasks.products_history(
	upload_id, article, article_supplier,
	brand, barcode, subject,
	price, upload_date, update_date,
	status, errorResponse) 
SELECT 
	b.upload_id, b.article, b.article_supplier,
	b.brand, b.barcode, b.subject,
	b.price, b.upload_date, b.update_date,
	b.status, b.errorResponse
FROM tasks.products_buffer AS b
WHERE upload_id = $1;`

	var res pgconn.CommandTag
	if res, err = _tx.Exec(ctx, query, uploadId); err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return fmt.Errorf("no rows were updated")
	}

	query = `DELETE FROM tasks.products_buffer WHERE upload_id = $1`
	if res, err = _tx.Exec(ctx, query, uploadId); err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return fmt.Errorf("no rows were updated")
	}
	if tx == nil {
		return _tx.Commit(ctx)
	}
	return nil
}

func (s *store) UpdateProductBufferErrorResponse(ctx context.Context, tx Transaction, products ...model.Product) (err error) {
	var _tx = tx
	if tx == nil {
		if _tx, err = s.pool.Begin(ctx); err != nil {
			return err
		}
		defer func() { _ = _tx.Rollback(ctx) }()
	}
	ctx, cancel := context.WithTimeout(ctx, s.cfg.Timeout)
	defer cancel()

	query :=
		`UPDATE tasks.products_buffer SET errorresponse = $1, update_date = $2 WHERE id = $3`

	for i := range products {
		row := _tx.QueryRow(ctx, query, products[i].ErrorResponse, time.Now().UTC(), products[i].ID)
		if err = row.Scan(); err != nil && !errors.Is(err, pgx.ErrNoRows) {
			return err
		}
	}
	if tx == nil {
		return _tx.Commit(ctx)
	}
	return nil
}
