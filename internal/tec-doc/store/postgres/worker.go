package postgres

import (
	"context"
	"fmt"
	"tec-doc/internal/tec-doc/model"
)

func (s *store) GetOldestTask(ctx context.Context, tx Transaction) (int64, error) {
	var (
		getOldestTaskQuery = `SELECT id FROM tasks WHERE status=$1 or status=$2 ORDER BY upload_date ASC LIMIT 1;`
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
	if err != nil {
		return 0, fmt.Errorf("can't exec getOldestTaskQuery: %w", err)
	}

	return t.ID, nil
}

func (s *store) GetProductsBufferWithStatus(ctx context.Context, tx Transaction, uploadID int64, limit int, offset int, status int) ([]model.Product, error) {
	var (
		getProductsBufferQuery = `SELECT id, upload_id, article, card_number, provider_article, manufacturer_article, brand, sku, category, price,
	upload_date, update_date, status, errorresponse FROM products_buffer WHERE upload_id = $1 and status=$2 LIMIT $3 OFFSET $4;`
		executor       Executor
		productsBuffer = make([]model.Product, 0)
	)

	executor = s.pool

	if tx != nil {
		executor = tx
	}

	rows, err := executor.Query(ctx, getProductsBufferQuery, uploadID, status, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("can't get products from buffer: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var p model.Product
		err := rows.Scan(&p.ID, &p.UploadID, &p.Article, &p.CardNumber, &p.ProviderArticle, &p.ManufacturerArticle, &p.Brand,
			&p.SKU, &p.Category, &p.Price, &p.UploadDate, &p.UpdateDate, &p.Status, &p.ErrorResponse)
		if err != nil {
			return nil, fmt.Errorf("can't get products from buffer: %w", err)
		}
		productsBuffer = append(productsBuffer, p)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("can't get products from buffer: %w", err)
	}

	return productsBuffer, nil
}

func (s *store) UpdateProductStatus(ctx context.Context, tx Transaction, productID int64, status int) error {

	var (
		updateProductStatusQuery = `UPDATE products_buffer SET status=$1 WHERE id=$2;`
		executor                 Executor
	)

	executor = s.pool
	if tx != nil {
		executor = tx
	}

	ctx, cancel := context.WithTimeout(ctx, s.cfg.Timeout)
	defer cancel()

	res, err := executor.Exec(ctx, updateProductStatusQuery, status, productID)
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
		updateTaskProductsNumberQuery = `UPDATE tasks SET products_failed=products_failed + $1, products_processed=products_processed + $2
                  WHERE id=$3;`
		executor Executor
	)

	executor = s.pool
	if tx != nil {
		executor = tx
	}

	ctx, cancel := context.WithTimeout(ctx, s.cfg.Timeout)
	defer cancel()

	res, err := executor.Exec(ctx, updateTaskProductsNumberQuery, productsFailed, productsProcessed, uploadID)
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
		updateTaskStatusQuery = `UPDATE task SET status=$1 WHERE id=$2;`
		executor              Executor
	)

	executor = s.pool
	if tx != nil {
		executor = tx
	}

	ctx, cancel := context.WithTimeout(ctx, s.cfg.Timeout)
	defer cancel()

	res, err := executor.Exec(ctx, updateTaskStatusQuery, status, uploadID)
	if err != nil {
		return fmt.Errorf("can't exec updateTaskStatusQuery: %w", err)
	}

	if res.RowsAffected() == 0 {
		return fmt.Errorf("no rows were updated")
	}

	return nil
}
