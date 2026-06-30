package service

type HealthStatus struct {
	Status string `json:"status"`
}

type HealthService struct{}

func NewHealthService() *HealthService {
	return &HealthService{}
}

func (s *HealthService) Status() HealthStatus {
	return HealthStatus{Status: "ok"}
}
