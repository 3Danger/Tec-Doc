package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"tec-doc/internal/config"
	"tec-doc/internal/model"
	"tec-doc/internal/store/postgres"
)

type Store interface {
	GetOldestTask(ctx context.Context, tx postgres.Transaction) (int64, int64, int64, error)
	GetProductFromBuffer(ctx context.Context, tx postgres.Transaction, uploadID, offset int64) (model.Product, error)
	UpdateProductStatus(ctx context.Context, tx postgres.Transaction, productID int64, status int) error
	UpdateTaskProductsNumber(ctx context.Context, tx postgres.Transaction, uploadID, productsFailed, productsProcessed int64) (bool, error)
	UpdateTaskStatus(ctx context.Context, tx postgres.Transaction, uploadID int64, status int) error
	Transaction(ctx context.Context) (postgres.Transaction, error)
}

func (s *store) Transaction(ctx context.Context) (postgres.Transaction, error) {
	return s.pool.Begin(ctx)
}

type store struct {
	cfg  *config.PostgresConfig
	pool *pgxpool.Pool
}

func NewStore(cfg *config.PostgresConfig) (*store, error) {
	pool, err := postgres.NewPool(cfg)
	if err != nil {
		return nil, fmt.Errorf("can't create pool: %w", err)
	}

	return &store{
		cfg:  cfg,
		pool: pool,
	}, nil
}

func (s *store) GetOldestTask(ctx context.Context, tx postgres.Transaction) (int64, int64, int64, error) {
	const getOldestTaskQuery = `SELECT id, products_processed + products_failed, products_total
	FROM tasks WHERE products_processed + products_failed < products_total ORDER BY upload_date DESC LIMIT 1;`

	var (
		executor postgres.Executor
		t        model.Task
		offset   int64
		total    int64
	)

	executor = s.pool
	if tx != nil {
		executor = tx
	}

	row := executor.QueryRow(ctx, getOldestTaskQuery)
	err := row.Scan(&t.ID, &offset, &total)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("can't exec getOldestTaskQuery: %w", err)
	}

	return t.ID, offset, total, nil
}

func (s *store) GetProductFromBuffer(ctx context.Context, tx postgres.Transaction, uploadID, offset int64) (model.Product, error) {
	const getProductFromBufferQuery = `SELECT id, upload_id, article, card_number, provider_article, manufacturer_article, brand, sku, category, price,
										upload_date, update_date, status, errorresponse FROM products_buffer WHERE upload_id = $1 LIMIT 1 OFFSET $2`

	var (
		executor postgres.Executor
		p        model.Product
	)

	executor = s.pool
	if tx != nil {
		executor = tx
	}

	row := executor.QueryRow(ctx, getProductFromBufferQuery, uploadID, offset)

	err := row.Scan(&p.ID, &p.UploadID, &p.Article, &p.CardNumber, &p.ProviderArticle, &p.ManufacturerArticle, &p.Brand,
		&p.SKU, &p.Category, &p.Price, &p.UploadDate, &p.UpdateDate, &p.Status, &p.ErrorResponse)

	if err != nil {
		return model.Product{}, fmt.Errorf("can't get oldest task: %w", err)
	}

	return p, nil
}

func (s *store) UpdateProductStatus(ctx context.Context, tx postgres.Transaction, productID int64, status int) error {
	const updateProductStatusQuery = `UPDATE products_buffer SET status=$1 WHERE id=$2;`

	var (
		executor postgres.Executor
	)

	executor = s.pool
	if tx != nil {
		executor = tx
	}
	res, err := executor.Exec(ctx, updateProductStatusQuery, status, productID)
	if err != nil {
		return fmt.Errorf("can't exec updateProductStatusQuery: %w", err)
	}

	if res.RowsAffected() == 0 {
		return fmt.Errorf("no rows were updated")
	}

	return nil
}

func (s *store) UpdateTaskProductsNumber(ctx context.Context, tx postgres.Transaction, uploadID, productsFailed, productsProcessed int64) (bool, error) {
	const updateTaskProductsNumberQuery = `UPDATE tasks SET products_failed=products_failed + $1, products_processed=products_processed + $2
                  WHERE id=$3 RETURNING products_processed, products_failed, products_total;`

	var (
		executor                  postgres.Executor
		proccessed, failed, total int64
	)

	executor = s.pool
	if tx != nil {
		executor = tx
	}
	row := executor.QueryRow(ctx, updateTaskProductsNumberQuery, productsFailed, productsProcessed, uploadID)
	err := row.Scan(&proccessed, &failed, &total)
	if err != nil {
		return false, fmt.Errorf("can't exec updateTaskProductsNumberQuery: %w", err)
	}

	return proccessed+failed == total, nil
}

func (s *store) UpdateTaskStatus(ctx context.Context, tx postgres.Transaction, uploadID int64, status int) error {
	const updateTaskStatusQuery = `UPDATE task SET status=$1 WHERE id=$2;`

	var (
		executor postgres.Executor
	)

	executor = s.pool
	if tx != nil {
		executor = tx
	}

	res, err := executor.Exec(ctx, updateTaskStatusQuery, status, uploadID)
	if err != nil {
		return fmt.Errorf("can't exec updateTaskStatusQuery: %w", err)
	}

	if res.RowsAffected() == 0 {
		return fmt.Errorf("no rows were updated")
	}

	return nil
}
