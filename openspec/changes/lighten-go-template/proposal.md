## Why

現在のリポジトリは「クリーンアーキテクチャ4層（domain/app/interfaces/infra）を `go test` で強制する」重厚なテンプレート骨組みだが、アプリ本体は未実装で、空ディレクトリ・厳密ルール・`旧サンプル名` 固有名・アーキテクチャテストだけが残っている。MVP・小規模チーム・AIエージェント実装には過剰で初速を落とす。`lightweight_template_implementation_guide.md` に従い、責務分離は保ちつつ `handler / service / repository` の3責務に整理した、すぐ使える軽量Goモノレポテンプレートへ作り替え、GitHub Template repository として配布する。

## What Changes

- **BREAKING**: 厳密4層アーキテクチャの強制を廃止。`backend/tests/arch/`（`structure_test.go` / `dependency_rule_test.go`）と4層前提の規約を削除する。
- backend を `handler / service / repository / config / db / router` の軽量構成へ作り替える（`backend/cmd/server/main.go` 起点）。
- 参照実装として `GET /health` を追加し、`{"status":"ok"}` を返す。`service` 層に最低1本のユニットテストを置く。
- frontend に React + TypeScript + Vite の最小スキャフォルド（`src/lib/api.ts`、`src/types/` 等）を追加する。
- `go.mod` のモジュールパスを `github.com/sudabon/<old-module>` から `旧サンプル名` 非依存の名前へリネームし、`旧サンプル名` 固有名を全廃する。
- `README.md` / `AGENTS.md` を軽量版向けに全面書き換え。`.claude` / `.codex` / `.cursor` の `rules/` から厳密4層強制部分を除去し、3責務ルールへ簡素化する。
- 配布用に `.env.example`、`docker-compose.yml`（PostgreSQLのみ）、`.gitignore` を整備する。
- GitHub リポジトリを **Template repository** に設定する。

## Capabilities

### New Capabilities
- `lightweight-go-backend`: `handler→service→repository` の3責務に整理したGo backend構成と、`GET /health` 参照実装、`service` のユニットテスト。
- `frontend-scaffold`: React + TypeScript + Vite の最小フロントエンド（APIクライアント集約・型定義配置の方針を含む）。
- `template-packaging`: 軽量版 README/AGENTS、エージェント向けルールの簡素化、`.env.example`・`docker-compose.yml`、重厚版資産（arch テスト・4層強制・`旧サンプル名` 固有名）の除去、GitHub Template repository 化。

### Modified Capabilities
<!-- openspec/specs/ は空のため、要件変更を伴う既存 capability はなし -->

## Impact

- **削除**: `backend/tests/arch/dependency_rule_test.go`、`backend/tests/arch/structure_test.go`、各 `rules/architecture.md` の4層強制部分、`.claude/rules/templates.md` の4層テンプレ、`旧サンプル名` 参照。
- **新規/書き換え**: `backend/cmd/server/main.go`、`backend/internal/{config,db,router,handler,service,repository}`、`backend/migrations/.gitkeep`、`frontend/`（Vite 一式）、`README.md`、`AGENTS.md`、`CLAUDE.md`、`.env.example`、`docker-compose.yml`、`.gitignore`、`go.mod`/`go.sum`。
- **依存**: Echo (v4) もしくは標準 `net/http`、PostgreSQL ドライバ（実DB接続は health では未使用でも可）。frontend は React/Vite/TypeScript。
- **配布**: `gh repo edit --template` による Template 化。フォーク利用者はモジュールパス・プロジェクト名のリネームのみで開始できる。
- **非互換**: 既存の重厚版を前提にしたフォーク手順・arch テストは無効化される。
