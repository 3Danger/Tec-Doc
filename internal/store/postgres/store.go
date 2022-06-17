package postgres

import (
	"context"
	"fmt"
	"tec-doc/internal/config"
	"tec-doc/internal/model"

	"github.com/jackc/pgx/v4/pgxpool"
)

//Store интерфейс описывающий методы для работы с БД
type Store interface {
	InsertContent(ctx *context.Context) error
	UpdateContent(ctx *context.Context) error
	GetContent(ctx *context.Context) (model.Content, error)
	DeleteContent(ctx *context.Context) error
}

type store struct {
	cfg  *config.Config
	pool *pgxpool.Pool
}

func NewStore(cfg *config.Config) (*store, error) {
	pool, err := NewPool(cfg)
	if err != nil {
		return nil, fmt.Errorf("can't create pool: %v", err)
	}

	return &store{
		cfg:  cfg,
		pool: pool,
	}, nil
}

func (s *store) InsertContent(ctx *context.Context) error {
	return nil
}

func (s *store) UpdateContent(ctx *context.Context) error {
	return nil
}

func (s *store) GetContent(ctx *context.Context) (model.Autopart, error) {
	return model.Autopart{}, nil
}

func (s *store) DeleteContent(ctx *context.Context) error {
	return nil
}

func NewPool(cfg *config.Config) (*pgxpool.Pool, error) {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?connect_timeout=%d", cfg.PostgresConfig.Username,
		cfg.PostgresConfig.Password, cfg.PostgresConfig.Host, cfg.PostgresConfig.Port,
		cfg.PostgresConfig.DbName, cfg.PostgresConfig.Timeout)

	connConf, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, err
	}

	connConf.MaxConns = cfg.PostgresConfig.MaxConns
	connConf.MinConns = cfg.PostgresConfig.MinConns
	pool, err := pgxpool.ConnectConfig(context.Background(), connConf)
	if err != nil {
		return nil, err
	}

	return pool, nil
}
