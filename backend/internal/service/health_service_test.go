package service

import (
	"context"
	"errors"
	"testing"
)

type fakeHealthPinger struct {
	err    error
	called bool
}

func (f *fakeHealthPinger) Ping(ctx context.Context) error {
	f.called = true
	return f.err
}

func TestHealthServiceStatus_通常時_OKを返す(t *testing.T) {
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

func TestDBHealthServiceCheck_Ping結果_DBステータスを返す(t *testing.T) {
	tests := []struct {
		name    string
		pingErr error
		want    DBHealthStatus
	}{
		{
			name: "ping succeeds",
			want: DBHealthStatus{Status: "ok"},
		},
		{
			name:    "ping fails",
			pingErr: errors.New("ping failed"),
			want:    DBHealthStatus{Status: "unavailable"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pinger := &fakeHealthPinger{err: tt.pingErr}
			service := NewDBHealthService(pinger)

			got := service.Check(context.Background())

			if got != tt.want {
				t.Fatalf("Check() = %#v, want %#v", got, tt.want)
			}
			if !pinger.called {
				t.Fatal("Check() did not call Ping()")
			}
		})
	}
}
