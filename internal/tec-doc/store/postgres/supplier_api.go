package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"tec-doc/pkg/model"
	"time"
)

func (s *store) CreateTask(ctx context.Context, tx Transaction, supplierID int64, userID int64, ip string, uploadDate time.Time) (int64, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	var (
		createTaskQuery = `INSERT INTO tasks.tasks (supplier_id, user_id, upload_date, update_date, IP, status, products_processed, products_failed, products_total)
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
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	var (
		getSupplierTaskHistoryQuery = `SELECT id, supplier_id, user_id, upload_date, update_date, status, products_processed, products_failed, products_total
								FROM tasks.tasks WHERE supplier_id = $1 ORDER BY upload_date LIMIT $2 OFFSET $3;`
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
		err = rows.Scan(&t.ID, &t.SupplierID, &t.UserID, &t.UploadDate,
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
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	var (
		getProductsFromHistoryQuery = `SELECT id, upload_id, article, article, manufacturer_article, brand, sku, category, price,
	upload_date, update_date, status, errorresponse FROM tasks.products_history WHERE upload_id = $1 LIMIT $2 OFFSET $3;`
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
		err = rows.Scan(&p.ID, &p.UploadID, &p.Article, &p.ArticleSupplier, &p.Subject,
			&p.Price, &p.UploadDate, &p.UpdateDate, &p.Status, &p.ErrorResponse)
		if err != nil {
			return nil, fmt.Errorf("can't get products from history: %w", err)
		}
		productsHistory = append(productsHistory, p)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("can't get products from history: %w", err)
	}

	return productsHistory, nil
}

func (s *store) SaveIntoBuffer(ctx context.Context, tx Transaction, products []model.Product) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	var (
		executor Executor
		rows     = make([][]interface{}, len(products))
	)

	executor = s.pool
	if tx != nil {
		executor = tx
	}

	for i, pr := range products {
		rows[i] = []interface{}{pr.UploadID, pr.Article, pr.ArticleSupplier, pr.Brand, pr.Barcode, pr.Price,
			time.Now().UTC(), time.Now().UTC(), pr.Status, pr.ErrorResponse}
	}

	copyCount, err := executor.CopyFrom(
		ctx,
		pgx.Identifier{"tasks", "products_buffer"},
		[]string{"upload_id", "article", "article_supplier", "brand", "barcode", "price", "upload_date", "update_date", "status", "errorresponse"},
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

func (s *store) GetProductsBuffer(ctx context.Context, tx Transaction, uploadID int64, limit int, offset int) ([]model.Product, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	var (
		getProductsBufferQuery = `SELECT id, upload_id, article, article_supplier, price,
		upload_date, update_date, status, errorresponse FROM tasks.products_buffer WHERE upload_id = $1 LIMIT $2 OFFSET $3;`
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
		err := rows.Scan(&p.ID, &p.UploadID, &p.Article, &p.ArticleSupplier, &p.Price, &p.UploadDate, &p.UpdateDate, &p.Status, &p.ErrorResponse)
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
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	var (
		executor Executor
		rows     = make([][]interface{}, len(products))
	)
	executor = s.pool
	if tx != nil {
		executor = tx
	}

	for i, pr := range products {
		rows[i] = []interface{}{pr.UploadID, pr.Article, pr.ArticleSupplier, pr.Price,
			pr.UploadDate, pr.UpdateDate, pr.Status, pr.ErrorResponse}
	}

	copyCount, err := executor.CopyFrom(
		ctx,
		pgx.Identifier{"tasks", "products_history"},
		[]string{"upload_id", "article", "article_supplier", "price", "upload_date", "update_date", "status", "errorresponse"},
		pgx.CopyFromRows(rows),
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
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	var (
		deleteFromBufferQuery = `DELETE FROM tasks.products_buffer WHERE upload_id=$1;`
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
