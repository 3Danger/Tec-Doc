package postgres

import (
	"context"
	"fmt"
	"tec-doc/internal/config"
	"tec-doc/internal/model"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	supplierStatusNew = iota
	supplierStatusProcess
	supplierStatusCompleted
	supplierStatusError
)

const (
	createSupplierQuery = `INSERT INTO suppliers (supplier_id, user_id, upload_date, updated_date, status)
							VALUES ($1, $2, $3, $4, $5)`
	saveIntoBufferQuery = `INSERT INTO products_buffer
							(task_id, article, brand, status, errorResponse, description)
							VALUES
    						($1, $2, $3, $4, $5);`
	moveFromBufferToHistoryQuery = `INSERT INTO products_history (id, task_id, article, brand, status, errorResponse, description)
									SELECT id, task_id, article, brand, status, errorResponse, description FROM products_buffer
									WHERE products_buffer.id NOT IN  (SELECT  id from products_history) 
											AND products_buffer.status = $1 AND products_buffer.task_id = $2;`
	deleteFromBufferQuery = `DELETE FROM products_buffer WHERE status = $1 AND products_buffer.task_id = $2;`
)

//Store интерфейс описывающий методы для работы с БД
type Store interface {
	CreateSupplier(ctx context.Context, supplierID int, userID int, ip string) error
	SaveIntoBuffer(ctx context.Context, uploadID int, article string, brand string, status int, errResponse string) error
	MoveFromBufferToHistory(ctx context.Context, status int, uploadID int) error
	InsertContent(ctx context.Context) error
	UpdateContent(ctx context.Context) error
	GetContent(ctx context.Context) (model.Content, error)
	DeleteContent(ctx context.Context) error
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

func (s *store) CreateSupplier(ctx context.Context, supplierID int, userID int, ip string) error {
	res, err := s.pool.Exec(ctx, createSupplierQuery, supplierID, userID,
		time.Now().UTC(), time.Now().UTC(), supplierStatusNew)

	if err != nil {
		return fmt.Errorf("can't exec createSupplier query: %v", err)
	}

	if res.RowsAffected() == 0 {
		return fmt.Errorf("no new suppliers created")
	}

	return nil
}

func (s *store) SaveIntoBuffer(ctx context.Context, uploadID int, article string, brand string, status int, errResponse string) error {
	res, err := s.pool.Exec(ctx, saveIntoBufferQuery, uploadID, article,
		brand, status, errResponse)

	if err != nil {
		return fmt.Errorf("can't exec saveIntoBuffer query: %v", err)
	}

	if res.RowsAffected() == 0 {
		return fmt.Errorf("no products saved into buffer table")
	}

	return nil
}

func (s *store) MoveFromBufferToHistory(ctx context.Context, status int, uploadID int) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("error initialising transaction: %v", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	res, err := s.pool.Exec(ctx, moveFromBufferToHistoryQuery, status, uploadID)
	if err != nil {
		return fmt.Errorf("can't exec MoveFromBufferToHistory query: %v", err)
	}

	if res.RowsAffected() == 0 {
		return fmt.Errorf("no products moved from buffer to history ")
	}

	res, err = s.pool.Exec(ctx, deleteFromBufferQuery, status, uploadID)
	if err != nil {
		return fmt.Errorf("can't exec deleteFromBuffer query: %v", err)
	}

	if res.RowsAffected() == 0 {
		return fmt.Errorf("no products deleted from buffer")
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("can't commit transaction: %v", err)
	}

	return nil
}

func (s *store) InsertContent(ctx context.Context) error {
	return nil
}

func (s *store) UpdateContent(ctx context.Context) error {
	return nil
}

func (s *store) GetContent(ctx context.Context) (*model.Article, error) {
	return nil, nil
}

func (s *store) DeleteContent(ctx context.Context) error {
	return nil
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
