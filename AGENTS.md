# AGENTS.md

このファイルはClaude Codeがリポジトリを操作する際のガイドラインである。
詳細なルールは `.claude/rules/` に分割して配置している。

## プロジェクト概要

モノレポ構成のWebアプリケーション。

- **バックエンド** (`backend/`): Go 1.23 / Echo。Clean Architectureに基づく4層構成。
- **フロントエンド** (`frontend/`): React / TypeScript / Vite / TailwindCSS。

## ディレクトリ構成

```
├── backend/                       # バックエンド
│   ├── cmd/
│   │   └── server/
│   │       └── main.go            # サーバプロセス管理（エントリポイント）
│   ├── internal/
│   │   └── todo_app/              # ソースルート（フォーク時にリネーム）
│   │       ├── app/               # Application層
│   │       │   ├── common/        # 共通の結果型・エラー型
│   │       │   ├── dto/           # Data Transfer Object（層間のデータ受け渡し）
│   │       │   ├── repository/    # データ永続化のインターフェース定義
│   │       │   ├── port/          # 外部サービス機能のインターフェース定義
│   │       │   └── usecase/       # ビジネスロジックの実行フロー
│   │       ├── domain/            # Domain層
│   │       │   ├── entity/        # エンティティ（ドメインモデル）
│   │       │   └── service/       # ドメインサービス（エンティティ横断ロジック）
│   │       ├── infra/             # Infrastructure層
│   │       │   ├── cli/           # CLIコマンド
│   │       │   ├── config/        # 設定（環境変数）
│   │       │   ├── notification/  # 通知（メール・Slack等）
│   │       │   ├── logging/       # ロギング（slog）
│   │       │   ├── persistence/   # データ永続化（GORMモデル・リポジトリ実装）
│   │       │   └── web/           # Web関連（Echoルーター・ミドルウェア）
│   │       └── interfaces/        # Interface層
│   │           ├── controller/    # 入力リクエスト処理
│   │           ├── presenter/     # 出力変換（domain固有型をviewmodel固有型へ）
│   │           └── viewmodel/     # レスポンス整形用ViewModel
│   ├── migrations/                # golang-migrate のSQLファイル
│   └── tests/
│       └── arch/                  # クリーンアーキテクチャのテスト
├── frontend/                      # フロントエンド
│   ├── src/
│   │   ├── components/            # UIコンポーネント（Shadcn/ui ベース）
│   │   ├── pages/                 # ページコンポーネント
│   │   ├── hooks/                 # カスタムフック
│   │   ├── lib/                   # ユーティリティ・API クライアント
│   │   ├── types/                 # 型定義
│   │   └── routes/                # React Router ルート定義
│   ├── e2e/                       # Playwright E2Eテスト
│   ├── package.json
│   └── vite.config.ts
├── go.mod
```

## 技術スタック

### バックエンド

- **言語**: Go 1.23+
- **フレームワーク**: Echo (v4)
- **パッケージ管理**: go modules
- **ORM**: GORM v2
- **DB**: PostgreSQL
- **マイグレーション**: golang-migrate
- **テスト**: go test（table-driven）
- **リンター**: golangci-lint / gofmt
- **静的解析**: go vet
- **DI**: コンストラクタ注入（手動配線。必要に応じ google/wire）

### フロントエンド

- **言語**: TypeScript
- **フレームワーク**: React
- **ビルドツール**: Vite
- **パッケージ管理**: pnpm
- **UIライブラリ**: Shadcn/ui + TailwindCSS
- **データ取得**: TanStack Query
- **ルーティング**: React Router
- **テスト**: Vitest（ユニット）/ Playwright（E2E）
- **リンター**: ESLint

## コマンド

### バックエンド

```bash
# 開発環境
go mod download                # 依存パッケージのダウンロード
go mod tidy                    # 依存の整理
go run ./backend/cmd/server    # 開発サーバー起動

# テスト
go test ./...                  # テスト全体実行
go test ./... -run TestName    # 特定テストの実行
go test ./... -failfast        # 最初の失敗で停止
go test ./backend/tests/arch/  # アーキテクチャテストのみ

# リント・静的解析・フォーマット
golangci-lint run              # リントチェック
go vet ./...                   # 静的解析
gofmt -l .                     # 未整形ファイルの一覧
gofmt -w .                     # フォーマット適用

# マイグレーション（golang-migrate）
migrate -path backend/migrations -database "$DATABASE_URL" up    # 最新まで適用
migrate -path backend/migrations -database "$DATABASE_URL" down 1 # 1つ前にロールバック
migrate create -ext sql -dir backend/migrations -seq <name>      # マイグレーション作成
```

### フロントエンド

```bash
cd frontend
pnpm install                   # 依存パッケージのインストール
pnpm dev                       # 開発サーバー起動
pnpm test                      # Vitest ユニットテスト実行
pnpm test:e2e                  # Playwright E2Eテスト実行
pnpm lint                      # ESLint チェック
pnpm build                     # プロダクションビルド
```

## 注意事項

- DBのURL等はすべて環境変数から取得する（`infra/config/` に集約）。
- バックエンド・フロントエンド間の型定義の乖離に注意する。APIスキーマを信頼の源泉とする。
- `interface` はGoの予約語のため、Interface層のパッケージ名は `interfaces`（複数形）とする。
