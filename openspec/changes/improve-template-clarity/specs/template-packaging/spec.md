## ADDED Requirements

### Requirement: README 手順の実効性

README の起動手順は、新規クローンから追加の暗黙知なしになぞれば動作しなければならない（SHALL）。効果のない手順を含めてはならない（MUST NOT）。Go アプリは `.env` を自動では読まないため、起動手順は環境変数を明示する形式とし、`.env.example` は「必要な環境変数のリファレンス」であることを明記する。

#### Scenario: 手順に効果のないステップがない

- **WHEN** README の起動手順を確認する
- **THEN** `cp .env.example .env` を前提とした手順がなく、`DATABASE_URL` 等を明示してサーバーを起動する手順が記載されている

#### Scenario: 新規クローンから動作を再現できる

- **WHEN** リポジトリを新規クローンし、README の手順を上から順に実行する
- **THEN** ブラウザでフロントエンドが「ok」を表示し、`/health/db` が DB 起動時 200・停止時 503 を返す

### Requirement: 環境変数ドキュメントの完全性

README の環境変数表は、`config` パッケージが読み込むすべての環境変数（`APP_ENV`・`APP_PORT`・`DATABASE_URL`・`CORS_ALLOW_ORIGINS`）と frontend の `VITE_API_BASE_URL` を、デフォルト値・用途とともに記載しなければならない（SHALL）。`.env.example` も同じ変数一覧を保つ。

#### Scenario: config と README が一致する

- **WHEN** `internal/config/config.go` が読む環境変数と README の環境変数表・`.env.example` を比較する
- **THEN** 過不足がない

### Requirement: エージェントルールの責務記述の同期

`handler / service / repository / config / db / router / main` の責務記述は、AGENTS.md と `.claude/rules/`・`.codex/rules/`・`.cursor/rules/` で同一内容を保たなければならない（SHALL）。router の責務は「Echo の構成とルート登録」、依存の配線（composition root）は「`cmd/server/main.go`」と記述する。

#### Scenario: 4箇所の責務記述が一致する

- **WHEN** AGENTS.md と3箇所の rules（architecture.md および backend-conventions.md）を比較する
- **THEN** router / main の責務記述を含め、内容が一致している
