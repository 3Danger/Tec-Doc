package main

import (
	"context"
	"fmt"
	"tec-doc/internal/config"
	"tec-doc/internal/model"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	supplierStatusNew = iota
	supplierStatusProcess
	supplierStatusCompleted
	supplierStatusError
)

//Store интерфейс описывающий методы для работы с БД
type Store interface {
	CreateTask(ctx context.Context, supplierID int, userID int, ip string) (time.Time, error)
	SaveIntoBuffer(ctx context.Context, products []model.Product) error
	GetSupplierTaskHistory(ctx context.Context, supplierID int, limit int, offset int) ([]model.Task, error)
	MoveFromBufferToHistory(ctx context.Context, uploadID int) error
	DeleteFromBuffer(ctx context.Context, uploadID int) error
	GetProductsHistory(ctx context.Context, uploadID int) ([]model.Product, error)
}

type store struct {
	cfg  *config.PostgresConfig
	pool *pgxpool.Pool
}

func NewStore(cfg *config.PostgresConfig) (*store, error) {
	pool, err := NewPool(cfg)
	if err != nil {
		return nil, fmt.Errorf("can't create pool: %v", err)
	}

	return &store{
		cfg:  cfg,
		pool: pool,
	}, nil
}

func (s *store) CreateTask(ctx context.Context, supplierID int, userID int, ip string, uploadDate time.Time) (int64, error) {
	createTaskQuery := `INSERT INTO tasks (supplier_id, user_id, upload_date, updated_date, IP, status)
							VALUES ($1, $2, $3, $4, $5, $6) RETURNING id;`

	row := s.pool.QueryRow(ctx, createTaskQuery, supplierID, userID,
		uploadDate, uploadDate, ip, supplierStatusNew)

	var taskID int64
	if err := row.Scan(&taskID); err != nil {
		return 0, fmt.Errorf("can't create task:: %w", err)
	}

	return taskID, nil
}

func (s *store) SaveIntoBuffer(ctx context.Context, products []model.Product) error {
	rows := make([][]interface{}, len(products))
	for i, pr := range products {
		r := make([]interface{}, 0)
		r = append(r, pr.UploadID, pr.Article, pr.Brand, pr.SKU, pr.Category, pr.Price,
			time.Now().UTC(), time.Now().UTC(), pr.Status, pr.ErrorResponse)
		rows[i] = r
	}

	copyCount, err := s.pool.CopyFrom(
		ctx,
		pgx.Identifier{"products_buffer"},
		[]string{"upload_id", "article", "brand", "sku", "category", "price",
			"upload_date", "update_date", "status", "errorresponse"},
		pgx.CopyFromRows(rows),
	)

	if err != nil {
		return fmt.Errorf("can't save into buffer: %v", err)
	}

	if copyCount == 0 {
		return fmt.Errorf("no products saved into products_buffer")
	}

	return nil
}

func (s *store) GetSupplierTaskHistory(ctx context.Context, supplierID int, limit int, offset int) ([]model.Task, error) {
	getSupplierTaskHistoryQuery := `SELECT id, supplier_id, user_id, IP, upload_date, update_date, status, products_processed, products_failed, products_total
								FROM tasks WHERE supplier_id = $1 ORDER BY upload_date LIMIT $2 OFFSET $3;`
	rows, err := s.pool.Query(ctx, getSupplierTaskHistoryQuery, supplierID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("can't get supplier task history: %v", err)
	}
	defer rows.Close()

	taskHistory := make([]model.Task, 0)
	for rows.Next() {
		var t model.Task
		err := rows.Scan(&t.ID, &t.SupplierID, &t.UserID, &t.IP, &t.UploadDate,
			&t.UpdateDate, &t.Status, &t.ProductsProcessed, &t.ProductsFailed, &t.ProductsFailed)
		if err != nil {
			return nil, fmt.Errorf("can't get tasks from history: %w", err)
		}
		taskHistory = append(taskHistory, t)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("can't get tasks from history: %w", err)
	}

	return taskHistory, nil
}

func (s *store) GetProductsFromBuffer(ctx context.Context, uploadID int) ([]model.Product, error) {
	getProductsBufferQuery := `SELECT upload_id, article, brand, sku, category, price,
	upload_date, update_date, status, errorresponse FROM products_buffer WHERE upload_id = $1;`

	rows, err := s.pool.Query(ctx, getProductsBufferQuery, uploadID)
	if err != nil {
		return nil, fmt.Errorf("can't get products from buffer: %v", err)
	}
	defer rows.Close()

	productsBuffer := make([]model.Product, 0)
	for rows.Next() {
		var p model.Product
		err := rows.Scan(&p.ID, &p.UploadID, &p.Article, &p.Brand, &p.SKU, &p.Category, &p.Price,
			&p.UploadDate, &p.UpdateDate, &p.Status, &p.ErrorResponse)
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

func (s *store) SaveProductsToHistory(ctx context.Context, products []model.Product) error {
	rowsBuf := make([][]interface{}, len(products))
	for i, pr := range products {
		r := make([]interface{}, 0)
		r = append(r, pr.UploadID, pr.Article, pr.Brand, pr.UploadDate, pr.UpdateDate, pr.Status, pr.ErrorResponse)
		rowsBuf[i] = r
	}

	copyCount, err := s.pool.CopyFrom(
		ctx,
		pgx.Identifier{"products_history"},
		[]string{"upload_id", "article", "brand", "upload_date", "update_date", "status", "errorresponse"},
		pgx.CopyFromRows(rowsBuf),
	)

	if err != nil {
		return fmt.Errorf("can't save products into history: %v", err)
	}

	if copyCount == 0 {
		return fmt.Errorf("no products saved into history")
	}

	return nil
}

func (s *store) DeleteFromBuffer(ctx context.Context, uploadID int) error {
	deleteFromBufferQuery := `DELETE FROM products_buffer WHERE upload_id=$1;`
	res, err := s.pool.Exec(ctx, deleteFromBufferQuery, uploadID)

	if err != nil {
		return fmt.Errorf("can't delete from buffer: %v", err)
	}

	if res.RowsAffected() == 0 {
		return fmt.Errorf("no products deleted from buffer")
	}

	return nil
}

func (s *store) GetProductsHistory(ctx context.Context, uploadID int, limit int, offset int) ([]model.Product, error) {
	getProductsFromHistoryQuery := `SELECT upload_id, article, brand, sku, category, price,
	upload_date, update_date, status, errorresponse FROM products_history WHERE upload_id = $1 LIMIT $2 OFFSET $3;`

	rows, err := s.pool.Query(ctx, getProductsFromHistoryQuery, uploadID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("can't get products from history: %v", err)
	}
	defer rows.Close()

	productsHistory := make([]model.Product, 0)
	for rows.Next() {
		var p model.Product
		err := rows.Scan(&p.ID, &p.UploadID, &p.Article, &p.Brand, &p.SKU, &p.Category, &p.Price,
			&p.UploadDate, &p.UpdateDate, &p.Status, &p.ErrorResponse)
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

func NewPool(cfg *config.PostgresConfig) (*pgxpool.Pool, error) {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?connect_timeout=%d", cfg.Username,
		cfg.Password, cfg.Host, cfg.Port,
		cfg.DbName, cfg.Timeout)

	connConf, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, err
	}

	connConf.MaxConns = cfg.MaxConns
	connConf.MinConns = cfg.MinConns
	pool, err := pgxpool.ConnectConfig(context.Background(), connConf)
	if err != nil {
		return nil, err
	}

	return pool, nil
}
