## Context

現リポジトリ `sudabon/lightweight_monorepo_go_project_template` は重厚版 Clean Architecture テンプレートの骨組みのみで、`backend/tests/arch/`（4層強制テスト2本）・`go.mod`（`github.com/sudabon/<old-module>`）・各種 `rules/`・`AGENTS.md`/`README.md` だけが存在し、アプリ本体・frontend は未実装。`lightweight_template_implementation_guide.md` に沿って、責務分離は残しつつ `handler/service/repository` の3責務へ整理し、in-place で軽量化して GitHub Template repository 化する。スコープは Go 版のみ（Python 版は対象外）。

## Goals / Non-Goals

**Goals:**

- 重厚版資産（arch テスト・4層強制ルール・`旧サンプル名` 固有名）を除去する。
- backend を `cmd/server` + `internal/{config,db,router,handler,service,repository}` の軽量構成へ再構築する。
- `GET /health` を `handler → service` の流れで実装し、DB 未起動でも起動・応答できる状態にする。
- React + TypeScript + Vite の最小 frontend スキャフォルドを追加する。
- README/AGENTS を軽量版へ全面書き換えし、`.env.example`・`docker-compose.yml` を整備、GitHub Template 化する。

**Non-Goals:**

- Python 版テンプレート（別リポジトリ）の作成。
- 実テーブルの CRUD・マイグレーション実体・認証等の業務機能（`repository/` と `migrations/` は雛形のみ）。
- TanStack Query / Shadcn/ui の標準採用（任意導入とし、本変更では入れない）。
- 将来の Clean Architecture 化（拡張余地は残すが本変更では行わない）。

## Decisions

### D1: HTTP フレームワークは Echo (v4)

`net/http` も候補だが、既存の `rules/backend-conventions.md` が Echo を前提にしており、ルーティング・ミドルウェア・JSON バインドが標準装備で軽量テンプレートの初速に資するため Echo v4 を採用する。
- 代替案: 標準 `net/http`（Go 1.22 の `ServeMux` メソッドパターン）。依存ゼロだが、ミドルウェア/バインドを自作する必要があり、拡張時の定型コードが増える。

### D2: モジュールパスは実リポジトリ名に合わせる

`go.mod` を `github.com/sudabon/<old-module>` → `github.com/sudabon/lightweight_monorepo_go_project_template` にリネームする。リモート名と一致させることで `go get`/`go install` の整合を保つ（末尾ハイフンは Go モジュールパスとして有効）。フォーク利用者は README 手順で自分のパスへ置換する。
- 代替案: `旧サンプル名` のまま、または汎用名。前者は固有名排除に反し、後者はリモート不一致を生む。

### D3: health は DB 非依存・3責務分離の参照実装

`handler.HealthHandler` がリクエストを受け、`service.HealthService.Status()` を呼び、結果を JSON 化して返す。`service` は `{"status":"ok"}` 相当の値のみ返し、DB へアクセスしない。これにより DB 未起動でもサーバ起動・`/health` 応答が可能で、責務分離の最小例を兼ねる。
- 代替案: health で DB ping を行う。可用性確認には有用だが、テンプレートの「すぐ起動」体験と衝突するため採用しない（README に拡張例として記載）。

### D4: DB 接続は最小・遅延・任意

`internal/db/postgres.go` は `database/sql` + pgx stdlib ドライバ（`github.com/jackc/pgx/v5/stdlib`）で `*sql.DB` を生成する関数を提供する。`main.go` は `/health` 起動に DB 接続を必須としない（`DATABASE_URL` 利用は実テーブル追加時から）。ORM（GORM 等）は実テーブルが出るまで導入しない。
- 代替案: 起動時に必須接続。テンプレート利用の障壁になるため不採用。

### D5: config は標準ライブラリのみ

`internal/config/config.go` は `os.Getenv` ＋デフォルト値で `APP_ENV`/`APP_PORT`/`DATABASE_URL` を読み込む。viper 等の外部依存は入れない。環境変数アクセスはこのパッケージに集約する。

### D6: frontend は Vite react-ts を最小化して採用

`pnpm create vite`（react-ts）相当の最小構成を作成し、`src/{components,pages,hooks,lib,types}/` を用意。API 呼び出しは `src/lib/api.ts` に集約し、`/health` を叩く例を置く。レスポンス型は `src/types/` に定義。`pnpm-lock.yaml` をコミットする。

### D7: GitHub Template 化は API 経由

`gh api --method PATCH repos/sudabon/lightweight_monorepo_go_project_template -F is_template=true` で `is_template` を有効化する。push 完了後に実施し、`gh repo view --json isTemplate` で確認する。

### D8: ルール簡素化は3責務へ統一

`.claude` / `.codex` / `.cursor` の `rules/architecture.md`・`templates.md` から4層強制部分を除去し、handler/service/repository の責務と禁止事項（handler から DB 直操作禁止 等）へ書き換える。3ツールで内容を一致させ重複矛盾を避ける。

## Risks / Trade-offs

- **3責務へ緩めることで設計逸脱（handler に SQL 直書き等）を機械的に防げない** → AGENTS.md/rules に禁止事項を明記し、`service` 参照実装で型を示す。将来 lint/コードレビューで補強。
- **Echo 採用で依存が1つ増える** → 単一の確立されたライブラリであり、拡張時の定型削減で相殺。net/http への差し替えは router/handler 局所で可能。
- **末尾ハイフンのモジュールパス** → Go 仕様上は有効だが見栄えが悪い。テンプレート性質上フォーク時に置換される前提で許容。気になる場合はリポジトリ名自体のリネームを別途検討。
- **frontend の lock ファイル肥大** → 最小依存に絞り、TanStack Query/Shadcn を入れないことで抑制。
- **in-place 改変による履歴の断絶** → 重厚版を参照したい利用者向けに、削除内容を proposal/Impact と README の移行メモに残す。

## Migration Plan

1. ブランチ `lighten-go-template` を作成して作業（main 直 push しない）。
2. 削除: `backend/tests/arch/` 2ファイル、`rules` の4層強制部分、`旧サンプル名` 参照。
3. backend 再構築 → `go mod tidy` → `go test ./...` / `go run ./cmd/server` で `/health` を検証。
4. frontend 追加 → `pnpm install && pnpm build` を検証。
5. README/AGENTS/.env.example/docker-compose/.gitignore 整備。
6. コミット（日本語・関心事分割）→ push → PR。
7. マージ後 `gh api` で `is_template=true`、`gh repo view` で確認。
- ロールバック: PR マージ前なら破棄、後なら revert。Template 設定は API で false に戻せる。

## Open Questions

- リポジトリ名末尾の `-` は意図的か誤りか（モジュールパスの見栄えに影響。現状は実名に合わせる方針）。
- `migrations/` のツール（golang-migrate 等）は本変更では雛形のみとし、実導入は将来課題でよいか。
