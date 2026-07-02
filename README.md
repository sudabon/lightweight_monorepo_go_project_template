# lightweight_monorepo_go_project_template

Go + React の軽量モノレポテンプレートです。backend は Echo v4 を使い、`handler / service / repository` の3責務に絞っています。frontend は React + TypeScript + Vite の最小構成です。

## 概要

- `GET /health` は DB 非依存で `{"status":"ok"}` を返します。
- `GET /health/db` は DB 疎通可能なら `{"status":"ok"}`、疎通不可なら 503 と `{"status":"unavailable"}` を返します。
- 環境変数の読み込みは `backend/internal/config` に集約します。
- DB接続初期化は `backend/internal/db` に置きますが、初期状態の起動には必須ではありません。
- frontend の API 呼び出しは `frontend/src/lib/api.ts` に集約し、レスポンス型は `frontend/src/types/` に置きます。

## ディレクトリ構成

```text
.
├── backend/
│   ├── cmd/server/main.go
│   ├── go.mod
│   ├── go.sum
│   ├── internal/
│   │   ├── config/
│   │   ├── db/
│   │   ├── handler/
│   │   ├── repository/
│   │   ├── router/
│   │   └── service/
│   └── migrations/
├── frontend/
│   ├── index.html
│   ├── package.json
│   ├── vite.config.ts
│   └── src/
│       ├── App.tsx
│       ├── components/
│       ├── hooks/
│       ├── lib/
│       ├── pages/
│       └── types/
├── docker-compose.yml
└── .env.example
```

## 起動方法

`.env.example` は必要な環境変数のリファレンスです。Go アプリは `.env` を自動では読み込まないため、起動時に環境変数を明示します。

```bash
docker compose up -d db

cd backend
DATABASE_URL='postgres://postgres:postgres@localhost:5432/app?sslmode=disable' go run ./cmd/server
```

別ターミナルで確認します。

```bash
curl http://localhost:8080/health
curl http://localhost:8080/health/db
```

DB が起動していない場合や `DATABASE_URL` 未設定の場合もサーバーは起動し、`/health/db` は 503 と `{"status":"unavailable"}` を返します。`/health` は従来通り DB 非依存で 200 を返します。

frontend は必要に応じて API の向き先を指定します。

```bash
cd frontend
pnpm install
VITE_API_BASE_URL=http://localhost:8080 pnpm dev
```

## テスト

```bash
cd backend
go test ./...
go vet ./...
gofmt -l .
```

```bash
cd frontend
pnpm test
pnpm build
```

## テスト方針

- service の業務判断は unit test を書きます。
- handler は主要な HTTP ステータスとレスポンス JSON をテストします。
- repository は実DBが必要な場合のみ integration test として追加します。
- `/health` はDB非依存を維持します。

## 開発ルール

- `handler` は HTTP 入出力だけを扱い、`service` を呼びます。
- `service` は業務判断を扱い、永続化が必要な場合だけ `repository` を呼びます。
- `repository` は DB・SQL・ドライバ固有処理を扱います。
- `handler` から DB を直接触らないでください。
- `os.Getenv` は `backend/internal/config` 以外で使わないでください。
- 小規模なレスポンス型は `service` に置いてよいです。ただし、HTTP専用の表現や画面都合の型が増えた場合は `handler` 側または `dto` に分離してください。
- frontend コンポーネントから直接 `fetch` を増やさず、`src/lib` の API クライアントに集約してください。

## DB接続

`db.OpenPostgres` は接続プールを初期化するだけで、接続確認は行いません。起動時にDB必須とする場合は、呼び出し側で `PingContext` を実行してください。ただし、初期テンプレートでは `/health` をDB非依存にするため、起動時DB接続確認は行いません。

マイグレーションは `backend/migrations/` に置き、必要になった時点で golang-migrate などの実行ツールを導入してください。

## 環境変数

| 変数 | デフォルト | 用途 |
|------|------------|------|
| `APP_ENV` | `local` | backend の実行環境名 |
| `APP_PORT` | `8080` | backend の待受ポート |
| `DATABASE_URL` | 空 | PostgreSQL 接続URL。`/health/db` の疎通確認で使用 |
| `CORS_ALLOW_ORIGINS` | `http://localhost:5173` | backend が許可する frontend の Origin（カンマ区切り） |
| `VITE_API_BASE_URL` | `http://localhost:8080` | frontend から backend を呼ぶURL |

## テンプレート利用時

このリポジトリから新規プロジェクトを作る場合は、`backend/go.mod` の module path、README、package 名を自分のリポジトリ名に合わせて変更してください。
