## ADDED Requirements

### Requirement: DB 死活確認エンドポイント

backend は `GET /health/db` を提供し、handler→service→repository→db の全連結を実演する参照実装としなければならない（SHALL）。DB へ疎通できる場合はステータス200と `{"status":"ok"}` を、疎通できない場合はステータス503と `{"status":"unavailable"}` を返す。repository は `PingContext` による接続確認のみを行い、疎通可否の判断（タイムアウト含む）は service に置く。既存の `GET /health` は引き続き DB 非依存を維持しなければならない（MUST）。

#### Scenario: DB 起動中は ok を返す

- **WHEN** PostgreSQL を起動し、`DATABASE_URL` を明示してサーバーを起動し、`GET /health/db` をリクエストする
- **THEN** HTTPステータス200と JSON ボディ `{"status":"ok"}` が返る

#### Scenario: DB に疎通できなくても 503 で応答する

- **WHEN** `DATABASE_URL` が未設定、または PostgreSQL が停止した状態でサーバーを起動し、`GET /health/db` をリクエストする
- **THEN** サーバー自体は正常に起動し、HTTPステータス503と JSON ボディ `{"status":"unavailable"}` が返る

#### Scenario: /health は DB 非依存のまま

- **WHEN** PostgreSQL を起動せずに `GET /health` をリクエストする
- **THEN** HTTPステータス200と `{"status":"ok"}` が返る

#### Scenario: テストは実 DB なしで成功する

- **WHEN** 実 DB を起動せずに `cd backend && go test ./...` を実行する
- **THEN** すべてのテストが成功する（service はフェイクの `HealthPinger`、handler はモック service で検証される）

### Requirement: composition root は main.go に置く

依存の生成と配線（config → db → repository → service → handler）は `cmd/server/main.go` に置かなければならない（SHALL）。`router` は Echo インスタンスの構成（ミドルウェア登録）とルート登録のみを担い、生成済みハンドラを引数で受け取る。

#### Scenario: router が下位レイヤーを import しない

- **WHEN** `internal/router/router.go` の import を確認する
- **THEN** `config`・`db`・`repository`・`service` パッケージへの import が存在しない（handler と Echo 関連パッケージ（`echo/v4`・`echo/v4/middleware`）のみ）

#### Scenario: 依存の組み立てが main.go で読める

- **WHEN** `cmd/server/main.go` を確認する
- **THEN** config 読み込み、DB プール生成、repository・service・handler の生成、router への受け渡しが一箇所に記述されている

### Requirement: CORS 許可

backend はフロントエンド開発サーバーからのクロスオリジンリクエストを許可しなければならない（SHALL）。許可オリジンは環境変数 `CORS_ALLOW_ORIGINS`（カンマ区切り、デフォルト `http://localhost:5173`）として `config` に集約する。

#### Scenario: 許可オリジンからのリクエストに CORS ヘッダが付く

- **WHEN** `Origin: http://localhost:5173` ヘッダ付きで `GET /health` をリクエストする
- **THEN** レスポンスに `Access-Control-Allow-Origin` ヘッダが含まれ、ブラウザ上のフロントエンドがレスポンスを読み取れる

### Requirement: Go モジュールは backend/ に置く

`go.mod` / `go.sum` は `backend/` 直下に置き、module は `github.com/sudabon/lightweight_monorepo_go_project_template/backend` としなければならない（SHALL)。モジュール境界とドキュメントの操作単位（`cd backend`）を一致させる。

#### Scenario: モジュール境界がディレクトリ境界と一致する

- **WHEN** リポジトリルートと `backend/` を確認する
- **THEN** ルートに `go.mod` が存在せず、`backend/go.mod` の module 行が `github.com/sudabon/lightweight_monorepo_go_project_template/backend` である

#### Scenario: backend ディレクトリからビルド・テストできる

- **WHEN** `cd backend && go build ./... && go test ./...` を実行する
- **THEN** すべて成功する

### Requirement: 自己規約との整合

エクスポートされる型・関数には `// Xxx は…` 形式の doc comment を付与しなければならない（SHALL）。health のテストは環境変数に依存してはならず（MUST NOT）、環境変数の読み込み挙動は `config` パッケージのユニットテストで検証する。

#### Scenario: エクスポートシンボルに doc comment がある

- **WHEN** `backend/internal/` 配下のエクスポートされた型・関数・メソッドを確認する
- **THEN** すべてに `// Xxx は…` 形式の doc comment が付与されている

#### Scenario: 環境変数のテストが config に集約されている

- **WHEN** backend のテストコードを確認する
- **THEN** `t.Setenv` の使用は `internal/config` のテストのみで、health の handler / service テストに環境変数への言及がない
