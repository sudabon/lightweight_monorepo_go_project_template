# 実装指示書: テンプレートのわかりやすさ改善

## 背景と目的

本リポジトリは [monorepo_go_project_template](https://github.com/sudabon/monorepo_go_project_template)（domain/app/infra/interfaces の4層 + アーキテクチャテスト）の軽量版として、handler/service/repository の3責務構成を採用しています。

現状の最大の問題は **「ルールは文書化されているが、コードで実演されていない」** ことです。repository は README のみ、`db.OpenPostgres` と `cfg.DatabaseURL` はどこからも使われず、handler→service→repository→db の流れを示す実例が1本もありません。テンプレートは「コピーして書き方を学ぶ」ものなので、実例の欠如がそのまま「わかりやすさがない」につながっています。

本指示書は、軽量さを維持したまま（4層化はしない）、パターンの実演・動作の正しさ・自己規約との整合を回復するための改修指示です。

## 改善提案一覧

| # | 優先度 | 対象 | 問題 | 提案 | 理由・効果 |
|---|--------|------|------|------|------------|
| 1 | 高 | backend 全体 | 3責務パターンの実例がゼロ（repository 空、db/DatabaseURL 未使用） | DB を使う縦切りサンプル `GET /health/db` を1本追加 | 全レイヤーの接続・DI・テストの見本が揃い、テンプレートとして自己説明的になる |
| 2 | 高 | backend/frontend 連携 | CORS 未設定。Vite dev (5173) から backend (8080) への fetch はブラウザにブロックされ、UI が常に offline 表示になる | Echo に CORS ミドルウェアを追加（許可オリジンは config 経由） | README 通りに起動して動く状態にする。テンプレートの第一印象を修復 |
| 3 | 中 | go.mod / go.sum | Go モジュールがリポジトリルートにあり、全コマンドを `cd backend` で実行させるドキュメントと食い違う | go.mod/go.sum を `backend/` に移動 | モジュール境界とディレクトリ境界が一致し、モノレポ構成が素直になる |
| 4 | 中 | router.go / main.go | 依存の生成が router に埋まっており、composition root が不明瞭 | 依存の組み立てを main.go に移し、router はルート登録専任にする | 依存の方向が main.go を読むだけで分かる。DB 導入時に router が肥大化しない |
| 5 | 中 | backend 全 .go ファイル | エクスポートシンボルに doc comment がなく、自らの規約（backend-conventions.md）に違反 | 全エクスポートに `// Xxx は…` 形式の doc comment を追加 | テンプレートが自分のルールの見本になる |
| 6 | 中 | 両テストファイル | health は DB 非依存なのに `t.Setenv("DATABASE_URL", …)` があり誤解を招く。service テストに重複アサート | 不要な Setenv と重複アサートを削除 | テストは真っ先にコピーされるので、ノイズの伝播を防ぐ |
| 7 | 中 | frontend | テスト基盤ゼロで testing.md（Vitest/Playwright 方針）と乖離。空の hooks/ の使い方も未提示 | データ取得ロジックを `useHealth` フックに抽出し、Vitest + テスト1本を追加 | hooks の配置・テストの書き方の見本が揃う |
| 8 | 低 | frontend/src/lib/api.ts | `as Promise<HealthResponse>` は型検証なしで、「unknown + 型ガード」規約の精神に反する | 型ガード `isHealthResponse` を追加 | 型ガードの見本を提供 |
| 9 | 低 | main.go / router.go | Recover ミドルウェアなし（panic で接続が落ちる）、graceful shutdown なし | `middleware.Recover()` とシグナルハンドリング + `e.Shutdown` を追加 | サーバーの基本作法の見本を提供 |
| 10 | 低 | .claude/.codex/.cursor rules | ルール文書が3箇所に手動コピーされておりドリフトのリスク | 同期チェックスクリプト（diff 比較）を追加 | 「同じ内容で配置」の保証を仕組み化 |
| 11 | 低 | migrations/ | マイグレーションの運用方法が未記載 | README に想定ツール（例: golang-migrate）を1行明記、または「必要時に導入」と記す | 空ディレクトリの意図を明確化 |
| 12 | 中 | README / config | README は `cp .env.example .env` を指示するが、config は `os.Getenv` のみで `.env` を誰も読まない（docker-compose も参照していない） | README の起動手順を「環境変数を明示して起動」に修正し、`.env.example` の役割を明記 | 手順通りにやっても効果のないステップを排除。#1 の動作確認の前提になる |

## 実施順序

1. **#4**（composition root）→ 2. **#1 + #12**（縦切りサンプルと .env 手順修正）→ 3. **#2**(CORS) → 4. **#5, #6**（規約整合）→ 5. **#7, #8**(frontend) → 6. **#3**（go.mod 移動。import 文字列が変わらないため独立しており、いつ実施してもよい）→ 7. **#9〜#11**（任意）

#3 は当初「構造変更を先に」と考えられますが、現状でも `cd backend` からのコマンドは親の go.mod を解決できて動作影響がないため、動くものを直す #1/#2 を優先します。

各ステップは独立した PR（1コミット1関心事）とし、コミットメッセージは日本語 + プレフィックス（`refactor:`, `feat:`, `fix:`, `test:`, `docs:`）を使ってください。

---

## 詳細指示

### #3 go.mod / go.sum を backend/ に移動

**現状**: go.mod がルートにあり module は `github.com/sudabon/lightweight_monorepo_go_project_template`。import パスは `…/backend/internal/…`。ドキュメントは全コマンドを `cd backend` で実行させるため、モジュールの実体と操作単位が食い違っています。

**変更内容**:

1. `git mv go.mod go.sum backend/`
2. `backend/go.mod` の module 行を `github.com/sudabon/lightweight_monorepo_go_project_template/backend` に変更。
   - これにより既存の import パス文字列（`…/backend/internal/…`）は**一切変更不要**です。
3. `cd backend && go mod tidy`
4. README の「ディレクトリ構成」「テンプレート利用時」、AGENTS.md の構成図を更新。更新対象は README・AGENTS.md・`.claude/.codex/.cursor` の rules のみとし、`openspec/changes/` 配下の過去の change 文書は履歴のため変更しない。

**受け入れ条件**: `cd backend && go build ./... && go test ./... && go vet ./...` が成功。ルートに go.mod が存在しない。

### #4 composition root を main.go に移す

**現状**: `backend/internal/router/router.go` の `New()` が service/handler を自前生成しています。DB が絡むと router が config・db にも依存し、依存関係が追いにくくなります。

**変更内容**:

1. `router.New` のシグネチャを、生成済みハンドラを受け取る形に変更:

```go
// New は Echo インスタンスを生成し、ルートを登録する。
func New(healthHandler *handler.HealthHandler) *echo.Echo {
	e := echo.New()
	e.GET("/health", healthHandler.Get)
	return e
}
```

2. `main.go` で config → service → handler → router の順に組み立てる（#1 実施後は db・repository もここに入る）。
3. ルール文書の更新（4箇所同期: `.claude/rules/architecture.md`、`.codex/`、`.cursor/`、AGENTS.md）:
   - 変更前: 「`router`: Echo ルート登録と依存の配線」
   - 変更後: 「`router`: Echo ルート登録。依存の配線（composition root）は `cmd/server/main.go`」

**受け入れ条件**: router パッケージが config・db・repository を import していない。既存テストが成功。

### #1 DB を使う縦切りサンプル `GET /health/db` を追加

**現状**: repository 層は README のみ。`db.OpenPostgres` と `cfg.DatabaseURL` は未使用。テンプレートの核心パターンを示すコードが存在しません。

**方針**: `/health` の DB 非依存は維持したまま（AGENTS.md の必須要件）、DB 死活確認用の `GET /health/db` を追加し、handler→service→repository→db の全連結を最小コードで実演します。

**変更内容**:

1. `backend/internal/repository/health_repository.go` を新規作成（templates.md の形式に従う）:

```go
// HealthRepository は DB の死活確認を行う repository。
type HealthRepository struct {
	db *sql.DB
}

func NewHealthRepository(db *sql.DB) *HealthRepository

// Ping は DB への接続確認を行う。
func (r *HealthRepository) Ping(ctx context.Context) error {
	return r.db.PingContext(ctx)
}
```

2. `backend/internal/service/db_health_service.go` を新規作成。repository へのインターフェースは利用側（service）に定義し、タイムアウトは業務判断として service に定数で持つ:

```go
// dbPingTimeout は DB 死活確認のタイムアウト。
const dbPingTimeout = 2 * time.Second

// HealthPinger は DB への接続確認を抽象化する。
type HealthPinger interface {
	Ping(ctx context.Context) error
}

// DBHealthStatus は DB 死活確認の結果。
type DBHealthStatus struct {
	Status string `json:"status"` // "ok" | "unavailable"
}

// DBHealthService は DB の死活を判断する service。
type DBHealthService struct {
	pinger HealthPinger
}

func NewDBHealthService(pinger HealthPinger) *DBHealthService

// Check は DB への疎通を確認し、状態を返す。
func (s *DBHealthService) Check(ctx context.Context) DBHealthStatus
```

3. handler は新規ファイルを作らず、既存の `HealthHandler` を拡張する（health という1リソースに1ハンドラ）。コンストラクタに `DBHealthService` 用のインターフェースを追加し、`GetDB` メソッドを追加する。service の結果が `ok` なら 200、`unavailable` なら 503 を返す。`echo.Context` から `c.Request().Context()` を取り出して service に渡す（context 伝播の見本）。
4. `main.go` で `db.OpenPostgres(cfg.DatabaseURL)` → `repository.NewHealthRepository` → `service.NewDBHealthService` → handler → router の順に配線。
   - `OpenPostgres` がエラーを返した場合は `log.Fatalf` で終了してよい（DSN 不正は設定ミスであり、「DB 未起動」とは別。DB 未起動は Ping 時に検出され 503 になるため、起動時 DB 非必須の要件は保たれる）。
   - `main.go` に `defer pool.Close()` を置く。
5. テスト追加:
   - service: `HealthPinger` のフェイク（成功/失敗）を使った table-driven test。`TestDBHealthService_Ping失敗_unavailableを返す` 形式の命名。
   - handler: 200 / 503 それぞれのステータスとレスポンス JSON を検証。
6. README・AGENTS.md に `/health/db` の説明を追記（「DB 未起動でもサーバーは起動し、`/health/db` は 503 を返す」）。

**受け入れ条件**:
- `docker compose up -d db` を実行し、`DATABASE_URL` を明示して起動（`.env` は自動では読まれない。#12 参照）:
  `DATABASE_URL='postgres://postgres:postgres@localhost:5432/app?sslmode=disable' go run ./cmd/server`
  → `curl localhost:8080/health/db` → `200 {"status":"ok"}`
- DB なし（`DATABASE_URL` 未設定または DB 停止中）: 同 → `503 {"status":"unavailable"}`。サーバー起動自体は成功する。
- `/health` は従来通り DB 非依存で 200。
- `go test ./...` は実 DB なしで成功する（repository の integration test は追加しない。テスト方針通り）。

### #2 CORS 対応

**現状**: Vite dev サーバー (http://localhost:5173) から backend (http://localhost:8080) への fetch はクロスオリジンですが、Echo はデフォルトで `Access-Control-Allow-Origin` を返しません。ブラウザがレスポンスを遮断するため、README の手順通りに起動しても UI は常に「offline」になります（テンプレートの動作確認で最初に踏む問題です）。

**変更内容**（推奨案: backend 側で CORS 許可）:

1. `config.Config` に `CORSAllowOrigins string` を追加（デフォルト `http://localhost:5173`、環境変数 `CORS_ALLOW_ORIGINS`、カンマ区切り）。
2. ミドルウェア登録は router の責務とする（Echo の構成は router に集約。#4 の「ルート登録専任」は「Echo インスタンスの構成 + ルート登録」と読み替え、rules 文書にもそう書く）。router が config パッケージに依存しないよう、値は引数で渡す:

```go
// New は Echo インスタンスを生成し、ミドルウェアとルートを登録する。
func New(allowOrigins []string, healthHandler *handler.HealthHandler) *echo.Echo {
	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: allowOrigins,
	}))
	e.GET("/health", healthHandler.Get)
	return e
}
```

3. `main.go` で `strings.Split(cfg.CORSAllowOrigins, ",")` を渡す。
4. `.env.example` と README の環境変数表に `CORS_ALLOW_ORIGINS` を追記。

**代替案**: Vite の `server.proxy` で `/health` を backend に転送する（Go 側変更ゼロ）。ただしエンドポイント追加のたびに proxy 設定が増えるため、テンプレートとしては backend 側 CORS を推奨します。

**受け入れ条件**: backend と `pnpm dev` を README 手順通りに起動し、ブラウザで UI が「ok」を表示する。

### #5 doc comment の追加

**現状**: backend-conventions.md は「公開する型・関数には doc comment を付与する（`// Xxx は…` 形式）」と定めていますが、`Config`, `Load`, `OpenPostgres`, `HealthStatus`, `HealthService`, `NewHealthService`, `Status`, `HealthHandler`, `NewHealthHandler`, `Get`, `router.New` のすべてに doc comment がありません。

**変更内容**: 上記全シンボル（および #1 で追加するシンボル）に `// Xxx は…` 形式の doc comment を付与します。

**受け入れ条件**: エクスポートシンボル全件に doc comment が付いている（`go vet` は doc comment の欠落を検出しないため、レビューで確認する。機械検査したい場合は `revive` の `exported` ルールを利用）。

### #6 テストのノイズ除去

**現状**:
- `health_handler_test.go` と `health_service_test.go` の冒頭に `t.Setenv("DATABASE_URL", "postgres://invalid:…")` がありますが、テスト対象は環境変数を一切読みません。「health のテストには DB 設定が要る」という誤解を植え付けます。
- `health_service_test.go` は `got.Status != "ok"` と `got != tt.want` の重複アサートがあります。

**変更内容**:
1. 両ファイルから `t.Setenv("DATABASE_URL", …)` を削除。
2. service テストのアサートを `got != tt.want` の1本に統一。
3. テスト関数名を testing.md の `TestXxx_条件_期待結果` 形式に合わせる（例: `TestHealthService_Status_okを返す`、`TestHealthHandler_Get_200とokを返す`）。
4. 代わりに `config/config_test.go` を新規追加し、環境変数のデフォルト値・上書きを table-driven で検証する。ここが `t.Setenv` の正しい使いどころの見本になる。

**受け入れ条件**: `go test ./...` 成功。health のテストに環境変数への言及がなく、config のテストにのみ `t.Setenv` がある。

### #7 useHealth フックの抽出と Vitest 導入

**現状**: `App.tsx` に状態機械・API 呼び出し・UI がすべて同居しています。`src/hooks/` は空（.gitkeep のみ）で、testing.md が定める Vitest の使い方を示すものが何もありません。

**変更内容**:

1. `frontend/src/hooks/useHealth.ts` を新規作成し、`LoadState` 型・`refresh`・`useEffect`（AbortController 含む）を移動。App.tsx は `const { state, refresh } = useHealth();` で消費する。
2. Vitest を導入: `pnpm add -D vitest jsdom @testing-library/react`。`package.json` に `"test": "vitest run"` を追加。
3. `vite.config.ts` に Vitest 設定を追加（これがないと hook テストは DOM 環境がなく失敗する）:

```ts
/// <reference types="vitest/config" />
export default defineConfig({
  plugins: [react()],
  test: { environment: "jsdom" },
});
```

4. `frontend/src/hooks/useHealth.test.ts` を作成。`vi.mock("../lib/api")` で `fetchHealth` をモックし、成功時 `ready`／失敗時 `error` になることを検証（testing.md 通り、テストは対象と同階層）。
5. README のテスト節に `pnpm test` を追記。

**受け入れ条件**: `pnpm test` と `pnpm build` が成功。App.tsx から useState/useEffect の import が消えている。

### #8 API レスポンスの型ガード

**変更内容**: `frontend/src/lib/api.ts` の `return response.json() as Promise<HealthResponse>;` を、`unknown` で受けて型ガードで絞り込む形に変更:

```ts
function isHealthResponse(value: unknown): value is HealthResponse {
  return (
    typeof value === "object" &&
    value !== null &&
    "status" in value &&
    typeof value.status === "string"
  );
}
```

（`in` 演算子の絞り込みを使えば `as` は不要。）不正な形状なら `throw new Error("Unexpected health response shape")`。

**受け入れ条件**: `pnpm build` 成功。レスポンスを `as HealthResponse` / `as Promise<HealthResponse>` で断定するコードが残っていない。

### #12 .env の扱いを実態に合わせる

**現状**: README の起動手順は `cp .env.example .env` から始まりますが、config は `os.Getenv` のみで `.env` ファイルを読みません。docker-compose.yml も値をハードコードしており `env_file` を参照していません。つまり `.env` を作っても誰も読まず、手順として無意味です（DB 非依存の `/health` では露見しませんが、#1 で DB を使い始めると「設定したのに繋がらない」事故になります）。

**変更内容**（dotenv ライブラリは導入しない。軽量さ優先）:

1. README の起動手順から `cp .env.example .env` を外し、環境変数を明示する形に変更:

```bash
docker compose up -d db
cd backend
DATABASE_URL='postgres://postgres:postgres@localhost:5432/app?sslmode=disable' go run ./cmd/server
```

2. `.env.example` の役割を README に1行明記: 「必要な環境変数の一覧（リファレンス）。Go アプリは `.env` を自動では読まない。direnv 等を使う場合は各自で読み込むこと」。

**受け入れ条件**: README の手順を新規クローンからなぞって、追加の暗黙知なしに `/health/db` の 200 が得られる。

### #9〜#11（任意・低優先度）

- **#9**: `e.Use(middleware.Recover())` を追加。main.go に `signal.NotifyContext` + `e.Shutdown(ctx)` の graceful shutdown を追加。
- **#10**: `.claude/rules/` を正とし、`.codex/rules/`・`.cursor/rules/` との diff を検査するシェルスクリプト（例: `scripts/check-rules-sync.sh`）を追加。
- **#11**: README の「DB接続」節に、マイグレーションツールは必要時に導入する旨と候補（golang-migrate 等)を1行追記。

---

## 全体の完了条件

```bash
cd backend && go test ./... && go vet ./... && gofmt -l .   # すべて成功・出力なし
cd frontend && pnpm test && pnpm build                       # すべて成功
```

加えて、README の手順を新規クローンからなぞって、ブラウザで「ok」表示・`/health/db` の 200/503 切り替えが確認できること。

## 注意事項

- `/health` の DB 非依存・起動時 DB 非必須は維持する（AGENTS.md の必須要件）。
- 4層（domain/app/infra/interfaces）への回帰はしない。本指示書の目的は「軽量なまま実演を揃える」こと。
- ルール文書を変更する場合は `.claude/` `.codex/` `.cursor/` の3箇所 + AGENTS.md + README を同時に更新する。
- このリポジトリは OpenSpec を採用しているため、実装に着手する場合は本指示書を元に `openspec/changes/` に change（proposal/design/tasks）を起票してから進めることを推奨する。過去の change 文書（`lighten-go-template`）は履歴であり遡って修正しない。
