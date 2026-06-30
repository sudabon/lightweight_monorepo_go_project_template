package service

import "testing"

func TestHealthServiceStatus(t *testing.T) {
	tests := []struct {
		name string
		want HealthStatus
	}{
		{
			name: "returns ok",
			want: HealthStatus{Status: "ok"},
		},
	}

	service := NewHealthService()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := service.Status()
			if got != tt.want {
				t.Fatalf("Status() = %#v, want %#v", got, tt.want)
			}
		})
	}
}
