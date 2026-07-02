## 1. composition root の移動（指示書 #4）

- [x] 1.1 作業ブランチ `improve-template-clarity` を作成する（main 直接変更しない）
- [x] 1.2 `router.New` を、生成済みハンドラを受け取るシグネチャへ変更する（router から service 生成コードを除去）
- [x] 1.3 `cmd/server/main.go` で config → service → handler → router の順に組み立てる
- [x] 1.4 rules 文書の「router: 依存の配線」を「main: 依存の配線（composition root）、router: Echo 構成とルート登録」へ更新する。対象は `.claude/`・`.codex/`・`.cursor/` 各 rules の `architecture.md` と `backend-conventions.md`（「Echo のルート登録と依存配線は router に置く」の記述）、および AGENTS.md
- [x] 1.5 `cd backend && go test ./... && go vet ./...` が成功し、router パッケージが config・db・repository を import していないことを確認する

## 2. DB 縦切りサンプルと .env 手順修正（指示書 #1, #12）

- [x] 2.1 `internal/repository/health_repository.go` を作成する（`*sql.DB` をコンストラクタで受け取り、`Ping(ctx)` で `PingContext` を呼ぶ）
- [x] 2.2 `internal/service/db_health_service.go` を作成する（`HealthPinger` インターフェースを service 側に定義、タイムアウト2秒を定数で持つ、`DBHealthStatus` を返す `Check(ctx)`）
- [x] 2.3 既存 `HealthHandler` に DB 用 service インターフェースと `GetDB` メソッドを追加する（ok なら 200、unavailable なら 503。`c.Request().Context()` を service へ渡す）
- [x] 2.4 `main.go` で `db.OpenPostgres(cfg.DatabaseURL)` → repository → service → handler を配線する（`OpenPostgres` エラー時は `log.Fatalf`、`defer pool.Close()`）
- [x] 2.5 router に `GET /health/db` を登録する
- [x] 2.6 service テストを追加する（`HealthPinger` のフェイクで成功/失敗の table-driven。`TestXxx_条件_期待結果` 命名）
- [x] 2.7 handler テストを追加する（200 / 503 のステータスとレスポンス JSON）
- [x] 2.8 README の起動手順から `cp .env.example .env` を外し、`DATABASE_URL='postgres://…' go run ./cmd/server` 形式へ修正。`.env.example` の役割（リファレンスであり自動では読まれない）を明記する
- [x] 2.9 README・AGENTS.md に `/health/db` を追記する（DB 未起動でもサーバーは起動し 503 を返す旨を含む）
- [x] 2.10 動作確認: DB あり → `/health/db` が 200 `{"status":"ok"}`、`DATABASE_URL` 未設定または DB 停止中 → 503 `{"status":"unavailable"}`（サーバー起動は成功）、`/health` は従来通り DB 非依存で 200。`go test ./...` は実 DB なしで成功する

## 3. CORS 対応（指示書 #2）

- [x] 3.1 `config.Config` に `CORSAllowOrigins` を追加する（環境変数 `CORS_ALLOW_ORIGINS`、デフォルト `http://localhost:5173`、カンマ区切り）
- [x] 3.2 router で `middleware.CORSWithConfig` を登録する（許可オリジンは `[]string` 引数で受け取り、main で `strings.Split` して渡す）
- [x] 3.3 `.env.example` と README の環境変数表に `CORS_ALLOW_ORIGINS` を追記する
- [x] 3.4 動作確認: backend と `pnpm dev` を README 手順通りに起動し、ブラウザで UI が「ok」を表示する

## 4. 自己規約との整合（指示書 #5, #6）

- [x] 4.1 全エクスポートシンボルに `// Xxx は…` 形式の doc comment を付与する（`go vet` では検出不可のためレビューで確認）
- [x] 4.2 health の handler / service テストから `t.Setenv("DATABASE_URL", …)` を削除し、service テストの重複アサートを `got != tt.want` に一本化する
- [x] 4.3 テスト関数名を `TestXxx_条件_期待結果` 形式へ揃える
- [x] 4.4 `internal/config/config_test.go` を追加する（デフォルト値・環境変数上書きの table-driven。`t.Setenv` の正しい使いどころの見本）

## 5. frontend の見本整備（指示書 #7, #8）

- [x] 5.1 `src/hooks/useHealth.ts` を作成し、`LoadState`・`refresh`・`useEffect`（AbortController 含む）を App.tsx から移動する
- [x] 5.2 Vitest を導入する（`pnpm add -D vitest jsdom @testing-library/react`、`package.json` に `"test": "vitest run"`、`vite.config.ts` に `test: { environment: "jsdom" }`）
- [x] 5.3 `src/hooks/useHealth.test.ts` を作成する（`vi.mock("../lib/api")` で成功時 `ready`／失敗時 `error` を検証）
- [x] 5.4 `src/lib/api.ts` の `as Promise<HealthResponse>` を `unknown` + `isHealthResponse` 型ガード（`in` 演算子の絞り込み）へ変更する
- [x] 5.5 README のテスト節に `pnpm test` を追記する
- [x] 5.6 `pnpm test && pnpm build` が成功することを確認する

## 6. go.mod の backend/ 移動（指示書 #3）

- [x] 6.1 `git mv go.mod go.sum backend/` を実行し、module 行を `github.com/sudabon/lightweight_monorepo_go_project_template/backend` へ変更する（import 文字列は変更不要）
- [x] 6.2 `cd backend && go mod tidy && go build ./... && go test ./...` が成功し、ルートに go.mod が存在しないことを確認する
- [x] 6.3 README（ディレクトリ構成・テンプレート利用時）と AGENTS.md の構成図を更新する（`openspec/changes/` の過去文書は履歴のため変更しない）

## 7. 任意（指示書 #9〜#11。スコープ外の推奨事項）

- [x] 7.1 （任意）`middleware.Recover()` と `signal.NotifyContext` + `e.Shutdown` による graceful shutdown を追加する
- [x] 7.2 （任意）`.claude/rules/` を正として `.codex/`・`.cursor/` との diff を検査する `scripts/check-rules-sync.sh` を追加する
- [x] 7.3 （任意）README に migrations の方針（必要時に golang-migrate 等を導入）を1行追記する

## 8. 仕上げ

- [x] 8.1 全体の完了条件を確認する: `cd backend && go test ./... && go vet ./... && gofmt -l .`（出力なし）、`cd frontend && pnpm test && pnpm build`
- [x] 8.2 新規クローンから README 手順をなぞり、ブラウザ「ok」表示と `/health/db` の 200/503 切り替えを確認する
- [x] 8.3 関心事ごとに日本語コミットメッセージで分割コミットする（プレフィックス: feat/fix/refactor/test/docs）
- [x] 8.4 push して PR を作成する（PR 本文に本 change `improve-template-clarity` と `docs/refactoring-instructions.md` へのリンクを含める）
