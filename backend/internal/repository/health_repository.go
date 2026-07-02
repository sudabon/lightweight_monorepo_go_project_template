package repository

import (
	"context"
	"database/sql"
)

// HealthRepository はDBの疎通確認を扱うrepositoryです。
type HealthRepository struct {
	pool *sql.DB
}

// NewHealthRepository はHealthRepositoryを生成します。
func NewHealthRepository(pool *sql.DB) *HealthRepository {
	return &HealthRepository{pool: pool}
}

// Ping はDB接続プールへ疎通確認を行います。
func (r *HealthRepository) Ping(ctx context.Context) error {
	return r.pool.PingContext(ctx)
}
