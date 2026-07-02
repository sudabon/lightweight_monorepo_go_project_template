package config

import "testing"

func TestLoad_環境変数_Defaultと上書きを返す(t *testing.T) {
	tests := []struct {
		name string
		env  map[string]string
		want Config
	}{
		{
			name: "defaults",
			want: Config{
				AppEnv:           "local",
				AppPort:          "8080",
				DatabaseURL:      "",
				CORSAllowOrigins: "http://localhost:5173",
			},
		},
		{
			name: "overrides",
			env: map[string]string{
				"APP_ENV":            "test",
				"APP_PORT":           "18080",
				"DATABASE_URL":       "postgres://example",
				"CORS_ALLOW_ORIGINS": "http://localhost:5173,http://localhost:6174",
			},
			want: Config{
				AppEnv:           "test",
				AppPort:          "18080",
				DatabaseURL:      "postgres://example",
				CORSAllowOrigins: "http://localhost:5173,http://localhost:6174",
			},
		},
	}

	keys := []string{"APP_ENV", "APP_PORT", "DATABASE_URL", "CORS_ALLOW_ORIGINS"}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, key := range keys {
				t.Setenv(key, "")
			}
			for key, value := range tt.env {
				t.Setenv(key, value)
			}

			got := Load()
			if got != tt.want {
				t.Fatalf("Load() = %#v, want %#v", got, tt.want)
			}
		})
	}
}
