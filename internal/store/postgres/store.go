package postgres

import (
	"context"
	"fmt"
	"runtime"
	"tec-doc/internal/config"
	"tec-doc/internal/model"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rs/zerolog"
)

type Store interface {
	InsertContent() error
	UpdateContent() error
	GetContent() (model.Content, error)
	DeleteContent() error
}

type store struct {
	ctx  context.Context
	cfg  *config.Config
	log  *zerolog.Logger
	pool *pgxpool.Pool
}

func NewStore(ctx context.Context, cfg *config.Config, log *zerolog.Logger) (*store, error) {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", cfg.PostgresConfig.Username,
		cfg.PostgresConfig.Password, cfg.PostgresConfig.Host, cfg.PostgresConfig.Port,
		cfg.PostgresConfig.DbName)
	pool, err := NewPool(connStr)
	if err != nil {
		return nil, fmt.Errorf("can't create pool: %v", err)
	}

	return &store{
		ctx:  ctx,
		cfg:  cfg,
		log:  log,
		pool: pool,
	}, nil
}

func (s *store) InsertContent() error {
	return nil
}

func (s *store) UpdateContent() error {
	return nil
}

func (s *store) GetContent() (model.Content, error) {
	return model.Content{}, nil
}

func (s *store) DeleteContent() error {
	return nil
}

func NewPool(connstr string) (*pgxpool.Pool, error) {
	connConf, err := pgxpool.ParseConfig(connstr)
	if err != nil {
		return nil, err
	}

	connConf.MaxConns = int32(runtime.NumCPU())
	pool, err := pgxpool.ConnectConfig(context.Background(), connConf)
	if err != nil {
		return nil, err
	}

	return pool, nil
}
