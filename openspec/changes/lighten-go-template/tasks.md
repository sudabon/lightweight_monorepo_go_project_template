## 1. 準備とクリーンアップ

- [x] 1.1 作業ブランチ `lighten-go-template` を作成する（main 直接変更しない）
- [x] 1.2 `backend/tests/arch/dependency_rule_test.go` と `backend/tests/arch/structure_test.go` を削除し、空になった `backend/tests/arch/` を整理する
- [x] 1.3 `go.mod` のモジュールパスを `github.com/sudabon/<old-module>` → `github.com/sudabon/lightweight_monorepo_go_project_template` にリネームする
- [x] 1.4 リポジトリ全体を `旧サンプル名` で検索し、ソース・設定・ドキュメントから固有名参照を除去する

## 2. backend 軽量構成の構築

- [x] 2.1 `backend/internal/config/config.go` を作成し、`os.Getenv` ＋デフォルト値で `APP_ENV`/`APP_PORT`/`DATABASE_URL` を読み込む（環境変数アクセスはここに集約）
- [x] 2.2 `backend/internal/db/postgres.go` を作成し、`database/sql` + pgx stdlib ドライバで `*sql.DB` を生成する関数を用意する（接続は遅延・任意、`/health` では未使用）
- [x] 2.3 `backend/internal/service/health_service.go` を作成し、`Status()` が `ok` 相当を返す（DB 非依存）
- [x] 2.4 `backend/internal/handler/health_handler.go` を作成し、`service` を呼んで `{"status":"ok"}` を JSON 返却する（SQL/業務判断を書かない）
- [x] 2.5 `backend/internal/router/router.go` を作成し、Echo で `GET /health` を登録する
- [x] 2.6 `backend/cmd/server/main.go` を作成し、config 読み込み → router 構築 → `APP_PORT` で起動する（DB 接続必須にしない）
- [x] 2.7 `backend/internal/repository/README.md` を作成し、実テーブル登場時に追加する旨を記載する（雛形のみ）
- [x] 2.8 `backend/migrations/.gitkeep` を作成する
- [x] 2.9 `cd backend && go mod tidy` で依存を解決する（Echo v4・pgx を取得）

## 3. backend のテストと動作検証

- [x] 3.1 `backend/internal/service/health_service_test.go` を table-driven で1本追加する
- [x] 3.2 `cd backend && go test ./...` が成功することを確認する（arch テストが存在しないこと含む）
- [x] 3.3 `cd backend && go run ./cmd/server` で起動し、`GET /health` が 200 と `{"status":"ok"}` を返すことを確認する
- [x] 3.4 `gofmt -l .` と `go vet ./...` を通す

## 4. frontend スキャフォルドの追加

- [x] 4.1 `frontend/` に Vite react-ts 最小構成（`package.json`/`vite.config.ts`/`index.html`/`src/main.tsx`/`src/App.tsx`）を作成する
- [x] 4.2 `src/{components,pages,hooks,lib,types}/` を作成する
- [x] 4.3 `src/lib/api.ts` に API 呼び出しを集約し、`/health` を叩く例を置く（コンポーネントから直接 fetch しない）
- [x] 4.4 `src/types/` に health レスポンス型を定義する
- [x] 4.5 `cd frontend && pnpm install && pnpm build` が成功することを確認し、`pnpm-lock.yaml` をコミット対象にする

## 5. ドキュメントとエージェントルールの整備

- [x] 5.1 `README.md` を軽量版へ全面書き換え（概要・ディレクトリ構成・起動方法・テスト・開発ルール・環境変数。`go run ./cmd/server` を含む）
- [x] 5.2 `AGENTS.md` を軽量版へ全面書き換え（backend 3責務ルール＋禁止事項、frontend ルール、主要コマンド）
- [x] 5.3 `CLAUDE.md` の参照内容が軽量版 AGENTS.md と整合することを確認する
- [x] 5.4 `.claude/rules/`・`.codex/rules/`・`.cursor/rules/` の `architecture.md`/`templates.md` から4層強制部分を除去し、handler/service/repository の3責務ルールへ統一する（3ツールで内容一致）

## 6. 配布用設定ファイル

- [x] 6.1 `.env.example` を作成（`APP_ENV=local`/`APP_PORT=8080`/`DATABASE_URL=postgres://...`）
- [x] 6.2 `docker-compose.yml` を作成（PostgreSQL のみ、余分なサービスを含めない）
- [x] 6.3 `.gitignore` を見直し、`.env` を除外していることを確認する

## 7. コミット・PR・Template 化

- [x] 7.1 関心事ごとに日本語コミットメッセージで分割コミットする（プレフィックス: feat/refactor/docs/chore）
- [ ] 7.2 push して PR を作成する
- [ ] 7.3 マージ後、`gh api --method PATCH repos/sudabon/lightweight_monorepo_go_project_template -F is_template=true` で Template 化する
- [ ] 7.4 `gh repo view --json isTemplate` で `isTemplate=true` を確認する
