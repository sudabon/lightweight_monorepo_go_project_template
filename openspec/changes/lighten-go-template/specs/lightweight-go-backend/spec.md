## ADDED Requirements

### Requirement: 3責務に整理した backend ディレクトリ構成

backend は厳密4層（domain/app/interfaces/infra）ではなく、`handler / service / repository` の3責務へ整理した構成を提供しなければならない（SHALL）。HTTP処理は `handler`、業務処理は `service`、DBアクセスは `repository` に配置する。`config`（環境変数読み込み）、`db`（DB接続初期化）、`router`（ルーティング定義）を補助パッケージとして持つ。

#### Scenario: 軽量構成のディレクトリが存在する

- **WHEN** リポジトリの `backend/` を確認する
- **THEN** `cmd/server/main.go`、`internal/config/`、`internal/db/`、`internal/router/`、`internal/handler/`、`internal/service/`、`internal/repository/`、`migrations/` が存在する

#### Scenario: 4層の空ディレクトリが残っていない

- **WHEN** `backend/internal/` 配下を確認する
- **THEN** `domain/`、`app/`、`interfaces/`、`infra/` のいずれのディレクトリも存在しない

### Requirement: 責務分離の禁止事項

HTTPハンドラに業務ロジックやSQLを書いてはならない（MUST NOT）。`handler` からDBを直接操作してはならず、DBアクセスは `repository` 経由で行う。`service` が必要に応じて `repository` を呼び出す。環境変数は `config` に集約し、各所で直接読み込んではならない。

#### Scenario: handler が DB を直接操作しない

- **WHEN** `internal/handler/` 配下のコードを確認する
- **THEN** SQL文字列や DB ドライバ／ORM の直接呼び出しが含まれず、`service` を呼び出している

#### Scenario: 環境変数の読み込みが config に集約される

- **WHEN** backend 全体で環境変数アクセスを確認する
- **THEN** 環境変数の読み込みは `internal/config` のみで行われている

### Requirement: health エンドポイント

backend は `GET /health` を提供し、ステータス200で `{"status":"ok"}` を返さなければならない（SHALL）。この処理は `handler → service` の流れで実装し、参照実装として責務分離を示す。

#### Scenario: health が ok を返す

- **WHEN** サーバ起動後に `GET /health` をリクエストする
- **THEN** HTTPステータス200と JSON ボディ `{"status":"ok"}` が返る

#### Scenario: サーバが起動する

- **WHEN** `cd backend && go mod tidy && go run ./cmd/server` を実行する
- **THEN** APIサーバが設定ポートで起動し、リクエストを受け付ける

### Requirement: service のユニットテスト

`service` 層には最低1本のユニットテストを置かなければならない（SHALL）。`go test ./...` が成功し、アーキテクチャテストはこの軽量版に含まれない。

#### Scenario: テストが成功する

- **WHEN** `cd backend && go test ./...` を実行する
- **THEN** すべてのテストが成功（PASS）し、`backend/tests/arch/` のテストは存在しない

### Requirement: 固有名の排除

backend のモジュールパス・パッケージ名・コードから `旧サンプル名` 固有名を排除しなければならない（MUST NOT 残存）。`go.mod` のモジュールパスは `旧サンプル名` を含まない名前にリネームする。

#### Scenario: 旧サンプル名 参照が残っていない

- **WHEN** リポジトリ全体を `旧サンプル名` で検索する
- **THEN** ソース・設定・ドキュメントのいずれにも `旧サンプル名` 参照が存在しない
