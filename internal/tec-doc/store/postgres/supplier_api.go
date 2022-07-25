package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"tec-doc/internal/tec-doc/model"
	"time"
)

func (s *store) CreateTask(ctx context.Context, tx Transaction, supplierID int64, userID int64, ip string, uploadDate time.Time) (int64, error) {
	var (
		createTaskQuery = `INSERT INTO tasks (supplier_id, user_id, upload_date, update_date, IP, status, products_processed, products_failed, products_total)
							VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id;`
		executor Executor
		taskID   int64
	)

	executor = s.pool
	if tx != nil {
		executor = tx
	}

	row := executor.QueryRow(ctx, createTaskQuery, supplierID, userID,
		uploadDate, uploadDate, ip, StatusNew, 0, 0, 0)

	if err := row.Scan(&taskID); err != nil {
		return 0, fmt.Errorf("can't create task:: %w", err)
	}

	return taskID, nil
}

func (s *store) GetSupplierTaskHistory(ctx context.Context, tx Transaction, supplierID int64, limit int, offset int) ([]model.Task, error) {
	var (
		getSupplierTaskHistoryQuery = `SELECT id, supplier_id, user_id, IP, upload_date, update_date, status, products_processed, products_failed, products_total
								FROM tasks WHERE supplier_id = $1 ORDER BY upload_date LIMIT $2 OFFSET $3;`
		executor    Executor
		taskHistory = make([]model.Task, 0)
	)

	executor = s.pool
	if tx != nil {
		executor = tx
	}

	rows, err := executor.Query(ctx, getSupplierTaskHistoryQuery, supplierID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("can't get supplier task history: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var t model.Task
		err = rows.Scan(&t.ID, &t.SupplierID, &t.UserID, &t.IP, &t.UploadDate,
			&t.UpdateDate, &t.Status, &t.ProductsProcessed, &t.ProductsFailed, &t.ProductsFailed)
		if err != nil {
			return nil, fmt.Errorf("can't get tasks from history: %w", err)
		}
		taskHistory = append(taskHistory, t)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("can't get tasks from history: %w", err)
	}

	return taskHistory, nil
}

func (s *store) GetProductsHistory(ctx context.Context, tx Transaction, uploadID int64, limit int, offset int) ([]model.Product, error) {
	var (
		getProductsFromHistoryQuery = `SELECT id, upload_id, article, card_number, provider_article, manufacturer_article, brand, sku, category, price,
	upload_date, update_date, status, errorresponse FROM products_history WHERE upload_id = $1 LIMIT $2 OFFSET $3;`
		executor        Executor
		productsHistory []model.Product
	)

	executor = s.pool
	if tx != nil {
		executor = tx
	}

	rows, err := executor.Query(ctx, getProductsFromHistoryQuery, uploadID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("can't get products from history: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var p model.Product
		err := rows.Scan(&p.ID, &p.UploadID, &p.Article, &p.CardNumber, &p.ProviderArticle, &p.ManufacturerArticle, &p.Brand,
			&p.SKU, &p.Category, &p.Price, &p.UploadDate, &p.UpdateDate, &p.Status, &p.ErrorResponse)
		if err != nil {
			return nil, fmt.Errorf("can't get products from history: %w", err)
		}
		productsHistory = append(productsHistory, p)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("can't get products from history: %w", err)
	}

	return productsHistory, nil
}

func (s *store) SaveIntoBuffer(ctx context.Context, tx Transaction, products []model.Product) error {
	var (
		executor Executor
		rows     = make([][]interface{}, len(products))
		r        = make([]interface{}, 0, len(products))
	)

	executor = s.pool
	if tx != nil {
		executor = tx
	}

	for i, pr := range products {
		r = append(r, pr.UploadID, pr.Article, pr.CardNumber, pr.ProviderArticle,
			pr.ManufacturerArticle, pr.Brand, pr.SKU, pr.Category, pr.Price,
			time.Now().UTC(), time.Now().UTC(), pr.Status, pr.ErrorResponse)
		rows[i] = r
		r = r[:]
	}

	copyCount, err := executor.CopyFrom(
		ctx,
		pgx.Identifier{"products_buffer"},
		[]string{"upload_id", "article", "card_number", "provider_article", "manufacturer_article",
			"brand", "sku", "category", "price", "upload_date", "update_date", "status", "errorresponse"},
		pgx.CopyFromRows(rows),
	)

	if err != nil {
		return fmt.Errorf("can't save into buffer: %w", err)
	}

	if copyCount == 0 {
		return fmt.Errorf("no products saved into products_buffer")
	}

	return nil
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
