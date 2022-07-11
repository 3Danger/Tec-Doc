package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgconn"
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
	CreateTask(ctx context.Context, tx Transaction, supplierID int64, userID int64, ip string, uploadDate time.Time) (int64, error)
	SaveIntoBuffer(ctx context.Context, tx Transaction, products []model.Product) error
	GetSupplierTaskHistory(ctx context.Context, tx Transaction, supplierID int64, limit int, offset int) ([]model.Task, error)
	GetProductsFromBuffer(ctx context.Context, tx Transaction, uploadID int64) ([]model.Product, error)
	SaveProductsToHistory(ctx context.Context, tx Transaction, products []model.Product) error
	DeleteFromBuffer(ctx context.Context, tx Transaction, uploadID int64) error
	GetProductsHistory(ctx context.Context, tx Transaction, uploadID int64, limit int, offset int) ([]model.Product, error)
	Transaction(ctx context.Context) (Transaction, error)
}

type Transaction interface {
	Executor
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}

func (s *store) Transaction(ctx context.Context) (Transaction, error) {
	return s.pool.Begin(ctx)
}

type Executor interface {
	Query(ctx context.Context, sql string, args ...interface{}) (rows pgx.Rows, err error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	CopyFrom(ctx context.Context, tableName pgx.Identifier, columnNames []string, rowSrc pgx.CopyFromSource) (int64, error)
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
	
	err = pool.Ping(context.Background())
	if err != nil {
		return nil, fmt.Errorf("can't connect pool: %v", err)
	}

	return &store{
		cfg:  cfg,
		pool: pool,
	}, nil
}

func (s *store) CreateTask(ctx context.Context, tx Transaction, supplierID int64, userID int64, ip string, uploadDate time.Time) (int64, error) {
	createTaskQuery := `INSERT INTO tasks (supplier_id, user_id, upload_date, update_date, IP, status, products_processed, products_failed, products_total)
							VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id;`
	var executor Executor
	executor = s.pool
	if tx != nil {
		executor = tx
	}

	row := executor.QueryRow(ctx, createTaskQuery, supplierID, userID,
		uploadDate, uploadDate, ip, supplierStatusNew, 0, 0, 0)

	var taskID int64
	if err := row.Scan(&taskID); err != nil {
		return 0, fmt.Errorf("can't create task:: %w", err)
	}

	return taskID, nil
}

func (s *store) SaveIntoBuffer(ctx context.Context, tx Transaction, products []model.Product) error {
	var executor Executor
	executor = s.pool
	if tx != nil {
		executor = tx
	}

	rows := make([][]interface{}, len(products))
	for i, pr := range products {
		r := make([]interface{}, 0)
		r = append(r, pr.UploadID, pr.Article, pr.CardNumber, pr.ProviderArticle,
			pr.ManufacturerArticle, pr.Brand, pr.SKU, pr.Category, pr.Price,
			time.Now().UTC(), time.Now().UTC(), pr.Status, pr.ErrorResponse)
		rows[i] = r
	}

	copyCount, err := executor.CopyFrom(
		ctx,
		pgx.Identifier{"products_buffer"},
		[]string{"upload_id", "article", "card_number", "provider_article", "manufacturer_article",
			"brand", "sku", "category", "price", "upload_date", "update_date", "status", "errorresponse"},
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

func (s *store) GetSupplierTaskHistory(ctx context.Context, tx Transaction, supplierID int64, limit int, offset int) ([]model.Task, error) {
	getSupplierTaskHistoryQuery := `SELECT id, supplier_id, user_id, IP, upload_date, update_date, status, products_processed, products_failed, products_total
								FROM tasks WHERE supplier_id = $1 ORDER BY upload_date LIMIT $2 OFFSET $3;`
	var executor Executor
	executor = s.pool
	if tx != nil {
		executor = tx
	}

	rows, err := executor.Query(ctx, getSupplierTaskHistoryQuery, supplierID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("can't get supplier task history: %v", err)
	}
	defer rows.Close()

	taskHistory := make([]model.Task, 0)
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

func (s *store) GetProductsFromBuffer(ctx context.Context, tx Transaction, uploadID int64) ([]model.Product, error) {
	getProductsBufferQuery := `SELECT id, upload_id, article, card_number, provider_article, manufacturer_article, brand, sku, category, price,
	upload_date, update_date, status, errorresponse FROM products_buffer WHERE upload_id = $1;`

	var executor Executor
	executor = s.pool
	if tx != nil {
		executor = tx
	}

	rows, err := executor.Query(ctx, getProductsBufferQuery, uploadID)
	if err != nil {
		return nil, fmt.Errorf("can't get products from buffer: %v", err)
	}
	defer rows.Close()

	productsBuffer := make([]model.Product, 0)
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
	var executor Executor
	executor = s.pool
	if tx != nil {
		executor = tx
	}

	rowsBuf := make([][]interface{}, len(products))
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
		return fmt.Errorf("can't save products into history: %v", err)
	}

	if copyCount == 0 {
		return fmt.Errorf("no products saved into history")
	}

	return nil
}

func (s *store) DeleteFromBuffer(ctx context.Context, tx Transaction, uploadID int64) error {
	deleteFromBufferQuery := `DELETE FROM products_buffer WHERE upload_id=$1;`
	var executor Executor
	executor = s.pool
	if tx != nil {
		executor = tx
	}

	res, err := executor.Exec(ctx, deleteFromBufferQuery, uploadID)

	if err != nil {
		return fmt.Errorf("can't delete from buffer: %v", err)
	}

	if res.RowsAffected() == 0 {
		return fmt.Errorf("no products deleted from buffer")
	}

	return nil
}

func (s *store) GetProductsHistory(ctx context.Context, tx Transaction, uploadID int64, limit int, offset int) ([]model.Product, error) {
	getProductsFromHistoryQuery := `SELECT id, upload_id, article, card_number, provider_article, manufacturer_article, brand, sku, category, price,
	upload_date, update_date, status, errorresponse FROM products_history WHERE upload_id = $1 LIMIT $2 OFFSET $3;`
	var executor Executor
	executor = s.pool
	if tx != nil {
		executor = tx
	}

	rows, err := executor.Query(ctx, getProductsFromHistoryQuery, uploadID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("can't get products from history: %v", err)
	}
	defer rows.Close()

	productsHistory := make([]model.Product, 0)
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

func NewPool(cfg *config.PostgresConfig) (*pgxpool.Pool, error) {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", cfg.Username,
		cfg.Password, cfg.Host, cfg.Port,
		cfg.DbName /*, cfg.Timeout*/)

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
