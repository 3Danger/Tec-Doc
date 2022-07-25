package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"tec-doc/internal/tec-doc/model"
)

func (s *store) GetProductsBuffer(ctx context.Context, tx Transaction, uploadID int64, limit int, offset int) ([]model.Product, error) {
	var (
		getProductsBufferQuery = `SELECT id, upload_id, article, card_number, provider_article, manufacturer_article, brand, sku, category, price,
	upload_date, update_date, status, errorresponse FROM products_buffer WHERE upload_id = $1 LIMIT $2 OFFSET $3;`
		executor       Executor
		productsBuffer = make([]model.Product, 0)
	)

	executor = s.pool

	if tx != nil {
		executor = tx
	}

	rows, err := executor.Query(ctx, getProductsBufferQuery, uploadID, limit, offset)
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

func (s *store) SaveProductsToHistory(ctx context.Context, tx Transaction, products []model.Product) error {
	var (
		executor Executor
		rowsBuf  = make([][]interface{}, len(products))
	)
	executor = s.pool
	if tx != nil {
		executor = tx
	}

	for i, pr := range products {
		r := make([]interface{}, 0)
		r = append(r, pr.UploadID, pr.Article, pr.CardNumber, pr.ProviderArticle,
			pr.ManufacturerArticle, pr.Brand, pr.SKU, pr.Category, pr.Price,
			pr.UploadDate, pr.UpdateDate, pr.Status, pr.ErrorResponse)
		rowsBuf[i] = r
	}

	copyCount, err := executor.CopyFrom(
		ctx,
		pgx.Identifier{"products_history"},
		[]string{"upload_id", "article", "card_number", "provider_article", "manufacturer_article",
			"brand", "sku", "category", "price", "upload_date", "update_date", "status", "errorresponse"},
		pgx.CopyFromRows(rowsBuf),
	)

	if err != nil {
		return fmt.Errorf("can't save products into history: %w", err)
	}

	if copyCount == 0 {
		return fmt.Errorf("no products saved into history")
	}

	return nil
}

func (s *store) DeleteFromBuffer(ctx context.Context, tx Transaction, uploadID int64) error {
	var (
		deleteFromBufferQuery = `DELETE FROM products_buffer WHERE upload_id=$1;`
		executor              Executor
	)

	executor = s.pool
	if tx != nil {
		executor = tx
	}

	res, err := executor.Exec(ctx, deleteFromBufferQuery, uploadID)

	if err != nil {
		return fmt.Errorf("can't delete from buffer: %w", err)
	}

	if res.RowsAffected() == 0 {
		return fmt.Errorf("no products deleted from buffer")
	}

	return nil
}

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
