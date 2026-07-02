package db

import (
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// OpenPostgres はPostgreSQL接続プールを初期化します。
func OpenPostgres(databaseURL string) (*sql.DB, error) {
	return sql.Open("pgx", databaseURL)
}
