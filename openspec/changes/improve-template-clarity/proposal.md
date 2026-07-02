## Why

軽量テンプレートの現状は「ルールは文書化されているが、コードで実演されていない」状態にある。repository は README のみ、`db.OpenPostgres` と `cfg.DatabaseURL` は未使用で、テンプレートの核心である handler→service→repository→db の流れを示す実例が1本もない。さらに CORS 未設定により README 手順通りに起動してもフロントエンドは常に offline 表示になり、README が指示する `cp .env.example .env` は誰にも読まれない。テンプレートは「コピーして書き方を学ぶ」ものであるため、これらがそのまま「わかりやすさがない（クリーンではない）」という評価につながっている。

詳細な分析と実装指示は `docs/refactoring-instructions.md` にある。本 change はそれを OpenSpec として正式化するもの。軽量さは維持し、4層（domain/app/infra/interfaces）へは回帰しない。

## What Changes

- DB を使う縦切りサンプル `GET /health/db` を追加し、handler→service→repository→db の接続・DI・テストの見本を揃える（`/health` の DB 非依存は維持）。
- CORS 対応を追加する（許可オリジンは config 経由、デフォルト `http://localhost:5173`）。README 手順通りの起動でフロントエンドが「ok」を表示できるようにする。
- 依存の組み立て（composition root）を router から `cmd/server/main.go` へ移し、router は Echo の構成とルート登録に専任させる。rules 文書の責務記述も更新する。
- README の起動手順を実態に合わせる（`.env` は自動では読まれないため、環境変数を明示する手順へ修正）。
- 全エクスポートシンボルに doc comment を付与し、自らの規約（backend-conventions.md）との矛盾を解消する。
- health テストから無意味な `t.Setenv("DATABASE_URL", …)` を除去し、代わりに `config.Load` の table-driven テストを追加する（`t.Setenv` の正しい使いどころの見本）。
- frontend のデータ取得ロジックを `src/hooks/useHealth.ts` に抽出し、Vitest を導入してテスト1本を追加する（testing.md の方針と実態の乖離解消）。API レスポンスは型ガードで検証する。
- **BREAKING**: `go.mod` / `go.sum` をリポジトリルートから `backend/` へ移動し、module を `github.com/sudabon/lightweight_monorepo_go_project_template/backend` にする（既存 import パス文字列は変更不要）。
- 任意（スコープ外の推奨事項としてタスクにのみ記載）: Recover ミドルウェア + graceful shutdown、rules 3箇所の同期チェックスクリプト、migrations 運用の README 追記。

## Capabilities

### New Capabilities

<!-- なし -->

### Modified Capabilities

- `lightweight-go-backend`: DB 縦切りサンプル（`/health/db`）、composition root の main.go 移動、CORS、doc comment、テストのノイズ除去、go.mod の backend/ 移動。
- `frontend-scaffold`: `useHealth` フック抽出、Vitest 導入、API レスポンスの型ガード。
- `template-packaging`: README 起動手順の実効性回復（`.env` の扱い明記）、環境変数表の更新、rules 文書の責務記述更新。

※ `openspec/specs/` は未アーカイブのため空。各 capability のデルタは `ADDED Requirements` として記述する（先行 change `lighten-go-template` と同じ扱い）。

## Impact

- **backend**: `cmd/server/main.go`、`internal/router/router.go`、`internal/config/config.go`（`CORS_ALLOW_ORIGINS` 追加）、`internal/handler/health_handler.go`（`GetDB` 追加）、`internal/service/db_health_service.go`（新規）、`internal/repository/health_repository.go`（新規）、各 `_test.go`、`go.mod`/`go.sum`（backend/ へ移動）。
- **frontend**: `src/hooks/useHealth.ts`（新規）、`src/hooks/useHealth.test.ts`（新規）、`src/App.tsx`、`src/lib/api.ts`、`package.json`、`vite.config.ts`（Vitest 設定）。
- **ドキュメント**: `README.md`、`AGENTS.md`、`.claude/rules/`・`.codex/rules/`・`.cursor/rules/`（router/main の責務記述。3箇所同期）、`.env.example`。
- **非互換**: module path 変更（`…/backend`）。テンプレート利用者のリネーム手順は README で更新する。`openspec/changes/lighten-go-template/` は履歴のため遡って修正しない。
- **トレーサビリティ**: 分析・指示の原本は `docs/refactoring-instructions.md`（#1〜#12 の項番を tasks から参照する）。
