# lightweight_monorepo_go_project_template

Go + React の軽量モノレポテンプレートです。backend は Echo v4 を使い、`handler / service / repository` の3責務に絞っています。frontend は React + TypeScript + Vite の最小構成です。

## 概要

- `GET /health` は DB 非依存で `{"status":"ok"}` を返します。
- 環境変数の読み込みは `backend/internal/config` に集約します。
- DB接続初期化は `backend/internal/db` に置きますが、初期状態の起動には必須ではありません。
- frontend の API 呼び出しは `frontend/src/lib/api.ts` に集約し、レスポンス型は `frontend/src/types/` に置きます。

## ディレクトリ構成

```text
.
├── backend/
│   ├── cmd/server/main.go
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
├── .env.example
└── go.mod
```

## 起動方法

```bash
cp .env.example .env
docker compose up -d db

cd backend
go run ./cmd/server
```

別ターミナルで確認します。

```bash
curl http://localhost:8080/health
```

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
pnpm build
```

## 開発ルール

- `handler` は HTTP 入出力だけを扱い、`service` を呼びます。
- `service` は業務判断を扱い、永続化が必要な場合だけ `repository` を呼びます。
- `repository` は DB・SQL・ドライバ固有処理を扱います。
- `handler` から DB を直接触らないでください。
- `os.Getenv` は `backend/internal/config` 以外で使わないでください。
- frontend コンポーネントから直接 `fetch` を増やさず、`src/lib` の API クライアントに集約してください。

## 環境変数

| 変数 | デフォルト | 用途 |
|------|------------|------|
| `APP_ENV` | `local` | backend の実行環境名 |
| `APP_PORT` | `8080` | backend の待受ポート |
| `DATABASE_URL` | 空 | PostgreSQL 接続URL。初期 `/health` では未使用 |
| `VITE_API_BASE_URL` | `http://localhost:8080` | frontend から backend を呼ぶURL |

## テンプレート利用時

このリポジトリから新規プロジェクトを作る場合は、`go.mod` の module path、README、package 名を自分のリポジトリ名に合わせて変更してください。
