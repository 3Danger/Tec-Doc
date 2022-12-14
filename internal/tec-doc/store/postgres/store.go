package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgconn"
	"tec-doc/internal/tec-doc/config"
	"tec-doc/pkg/model"
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

type Store interface {
	CreateTask(ctx context.Context, tx Transaction, supplierIdString string, supplierID, userID, productsTotal int64, ip string, uploadDate time.Time) (int64, error)
	SaveIntoBuffer(ctx context.Context, tx Transaction, products []model.Product) error
	GetSupplierTaskHistory(ctx context.Context, tx Transaction, supplierID int64, limit int, offset int) ([]model.Task, error)
	GetProductsBuffer(ctx context.Context, tx Transaction, uploadID int64, limit int, offset int) ([]model.Product, error)
	MoveProductsToHistoryByUploadId(ctx context.Context, tx Transaction, uploadId int64) (err error)
	MarkTaskAsCompletedAndMoveProductsToHistory(ctx context.Context, uploadID int64) error
	UpdateProductBuffer(ctx context.Context, tx Transaction, products *model.Product) (err error)
	SaveProductsToHistory(ctx context.Context, tx Transaction, products []model.Product) error
	DeleteFromBuffer(ctx context.Context, tx Transaction, uploadID int64) error
	GetProductsHistory(ctx context.Context, tx Transaction, uploadID int64, limit int, offset int) ([]model.Product, error)
	GetProductsHistoryWithStatus(ctx context.Context, tx Transaction, uploadID, status int64, limit int, offset int) ([]model.Product, error)

	GetOldestTask(ctx context.Context, tx Transaction) (uploadId int64, supplierIdStr string, err error)
	GetProductsBufferWithStatus(ctx context.Context, tx Transaction, uploadID int64, limit int, offset int, status int) ([]model.Product, error)
	UpdateTaskProductsNumber(ctx context.Context, tx Transaction, uploadID, productsFailed, productsProcessed int64) error
	UpdateTaskStatus(ctx context.Context, tx Transaction, uploadID int64, status int) error

	Transaction(ctx context.Context) (Transaction, error)
	Stop()
}

type Transaction interface {
	Executor
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}

func (s *store) Transaction(ctx context.Context) (Transaction, error) {
	return s.pool.Begin(ctx)
}

func (s *store) Stop() {
	s.pool.Close()
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

func NewStore(ctx context.Context, cfg *config.PostgresConfig) (Store, error) {
	pool, err := NewPool(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("can't create pool: %v", err)
	}

	err = pool.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("can't connect pool: %v", err)
	}

	return &store{
		cfg:  cfg,
		pool: pool,
	}, nil
}

func NewPool(ctx context.Context, cfg *config.PostgresConfig) (*pgxpool.Pool, error) {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.DbName)

	connConf, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, errors.New("failed to parse config")
	}

	connConf.MaxConns = cfg.MaxConns
	connConf.MinConns = cfg.MinConns
	pool, err := pgxpool.ConnectConfig(ctx, connConf)
	if err != nil {
		return nil, errors.New("failed to connect config")
	}
	return pool, nil
}
