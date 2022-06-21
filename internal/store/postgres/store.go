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
	supplierTaskHistoryQuery = `SELECT INTO suppliers (supplier_id, user_id, upload_date, updated_date, status)
							VALUES ($1, $2, $3, $4, $5)`
)

//Store интерфейс описывающий методы для работы с БД
type Store interface {
	CreateSupplier(ctx context.Context, supplierID int, userID int, ip string) error
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
		return fmt.Errorf("can't exec database query: %v", err)
	}

	if res.RowsAffected() == 0 {
		return fmt.Errorf("can't create new supplier")
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
