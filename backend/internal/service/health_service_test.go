package service

import "testing"

func TestHealthServiceStatus(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://invalid:invalid@127.0.0.1:1/invalid")

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
			if got.Status != "ok" {
				t.Fatalf("Status().Status = %q, want %q", got.Status, "ok")
			}
			if got != tt.want {
				t.Fatalf("Status() = %#v, want %#v", got, tt.want)
			}
		})
	}
}
