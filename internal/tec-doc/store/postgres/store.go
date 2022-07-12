package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgconn"
	"tec-doc/internal/tec-doc/config"
	"tec-doc/internal/tec-doc/model"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	StatusNew = iota
	StatusProcess
	StatusCompleted
	StatusError
)

//Store интерфейс описывающий методы для работы с БД
type Store interface {
	CreateTask(ctx context.Context, tx Transaction, supplierID int64, userID int64, ip string, uploadDate time.Time) (int64, error)
	SaveIntoBuffer(ctx context.Context, tx Transaction, products []model.Product) error
	GetSupplierTaskHistory(ctx context.Context, tx Transaction, supplierID int64, limit int, offset int) ([]model.Task, error)
	GetProductsBuffer(ctx context.Context, tx Transaction, uploadID int64, limit int, offset int) ([]model.Product, error)
	SaveProductsToHistory(ctx context.Context, tx Transaction, products []model.Product) error
	DeleteFromBuffer(ctx context.Context, tx Transaction, uploadID int64) error
	GetProductsHistory(ctx context.Context, tx Transaction, uploadID int64, limit int, offset int) ([]model.Product, error)

	GetOldestTask(ctx context.Context, tx Transaction) (int64, error)
	GetProductsBufferWithStatus(ctx context.Context, tx Transaction, uploadID int64, limit int, offset int, status int) ([]model.Product, error)
	UpdateProductStatus(ctx context.Context, tx Transaction, productID int64, status int) error
	UpdateTaskProductsNumber(ctx context.Context, tx Transaction, uploadID, productsFailed, productsProcessed int64) error
	UpdateTaskStatus(ctx context.Context, tx Transaction, uploadID int64, status int) error

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
		return nil, fmt.Errorf("can't create pool: %w", err)
	}

	return &store{
		cfg:  cfg,
		pool: pool,
	}, nil
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

func (s *store) SaveIntoBuffer(ctx context.Context, tx Transaction, products []model.Product) error {
	var (
		executor Executor
		rows     = make([][]interface{}, len(products))
	)

	executor = s.pool
	if tx != nil {
		executor = tx
	}

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
		return fmt.Errorf("can't save into buffer: %w", err)
	}

	if copyCount == 0 {
		return fmt.Errorf("no products saved into products_buffer")
	}

	return nil
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

func (s *store) GetProductsHistory(ctx context.Context, tx Transaction, uploadID int64, limit int, offset int) ([]model.Product, error) {
	var (
		getProductsFromHistoryQuery = `SELECT id, upload_id, article, card_number, provider_article, manufacturer_article, brand, sku, category, price,
	upload_date, update_date, status, errorresponse FROM products_history WHERE upload_id = $1 LIMIT $2 OFFSET $3;`
		executor        Executor
		productsHistory = make([]model.Product, 0)
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
