package service

import (
	"context"
	"time"
)

const dbHealthTimeout = 2 * time.Second

// HealthPinger はDB疎通確認に必要な振る舞いを表します。
type HealthPinger interface {
	Ping(ctx context.Context) error
}

// DBHealthStatus はDB疎通確認のレスポンスです。
type DBHealthStatus struct {
	Status string `json:"status"`
}

// DBHealthService はDB疎通確認の業務判断を扱います。
type DBHealthService struct {
	pinger HealthPinger
}

// NewDBHealthService はDBHealthServiceを生成します。
func NewDBHealthService(pinger HealthPinger) *DBHealthService {
	return &DBHealthService{pinger: pinger}
}

// Check はDB疎通確認を行い、結果ステータスを返します。
func (s *DBHealthService) Check(ctx context.Context) DBHealthStatus {
	ctx, cancel := context.WithTimeout(ctx, dbHealthTimeout)
	defer cancel()

	if s.pinger == nil {
		return DBHealthStatus{Status: "unavailable"}
	}
	if err := s.pinger.Ping(ctx); err != nil {
		return DBHealthStatus{Status: "unavailable"}
	}

	return DBHealthStatus{Status: "ok"}
}
