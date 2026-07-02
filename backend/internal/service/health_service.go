package service

// HealthStatus はヘルスチェックのレスポンスです。
type HealthStatus struct {
	Status string `json:"status"`
}

// HealthService はDB非依存のヘルスチェックを扱います。
type HealthService struct{}

// NewHealthService はHealthServiceを生成します。
func NewHealthService() *HealthService {
	return &HealthService{}
}

// Status はDB非依存のヘルスチェック結果を返します。
func (s *HealthService) Status() HealthStatus {
	return HealthStatus{Status: "ok"}
}
