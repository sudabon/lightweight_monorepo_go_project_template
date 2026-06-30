# AGENTS.md

このリポジトリは Go + React の軽量モノレポテンプレートです。詳細ルールは `.claude/rules/`、`.codex/rules/`、`.cursor/rules/` に同じ内容で配置しています。

## 構成

- backend: Go 1.23+ / Echo v4 / PostgreSQL 接続用 pgx
- frontend: React / TypeScript / Vite / pnpm

```text
backend/
  cmd/server/main.go
  internal/config/
  internal/db/
  internal/handler/
  internal/repository/
  internal/router/
  internal/service/
  migrations/
frontend/
  src/components/
  src/hooks/
  src/lib/
  src/pages/
  src/types/
```

## backend 3責務ルール

- `handler` は HTTP 入出力を扱い、`service` を呼ぶ。
- `service` は業務判断を扱い、永続化が必要な場合に `repository` を呼ぶ。
- `repository` は DB・SQL・ドライバ固有処理を扱う。
- `handler` から DB を直接触らない。
- SQL を `handler` に書かない。
- 環境変数アクセスは `backend/internal/config` に集約し、各所で `os.Getenv` を呼ばない。
- DB 接続初期化は `backend/internal/db` に置く。起動時に DB を必須にしない。
- 小規模なレスポンス型は `service` に置いてよい。ただし、HTTP専用の表現や画面都合の型が増えた場合は `handler` 側または `dto` に分離する。
- `db.OpenPostgres` は接続プールを初期化するだけで、接続確認は行わない。起動時にDB必須とする場合は、呼び出し側で `PingContext` を実行する。

## frontend ルール

- API 呼び出しは `frontend/src/lib/` に集約する。
- API レスポンス型は `frontend/src/types/` に置く。
- コンポーネントから直接 `fetch` を増やさない。
- TanStack Query、TailwindCSS、UIライブラリは必要になった時点で導入する。

## 主要コマンド

```bash
cd backend
go mod tidy
go run ./cmd/server
go test ./...
go vet ./...
gofmt -l .
```

```bash
cd frontend
pnpm install
pnpm dev
pnpm build
```

## テスト方針

- service の業務判断は unit test を書く。
- handler は主要なHTTPステータスとレスポンスJSONをテストする。
- repository は実DBが必要な場合のみ integration test として追加する。
- `/health` はDB非依存を維持する。

## 変更方針

- 既存の公開APIを変える場合は事前に確認する。
- 変更は最小限に保ち、既存スタイルに合わせる。
- テストとビルドが利用できる場合は実行する。
- `.env` はコミットしない。配布用の値は `.env.example` に置く。
