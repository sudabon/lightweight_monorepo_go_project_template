# Design

## Goals / Non-Goals

**Goals**

- handler→service→repository→db の全連結を、最小コードの実例（`/health/db`）としてテンプレート内に揃える。
- README 手順通りに起動すれば動く状態にする（CORS・`.env` の実効性）。
- テンプレート自身が自らの規約（doc comment、テスト方針、環境変数集約）の見本になる。

**Non-Goals**

- 4層（domain/app/infra/interfaces）やアーキテクチャテストへの回帰。
- CRUD サンプル・マイグレーションツール・dotenv ライブラリ・TanStack Query 等の新規依存追加。
- `openspec/changes/lighten-go-template/`（履歴）の遡及修正。

## 前提

- 軽量テンプレートの思想（3責務 + 補助パッケージ、起動時 DB 非必須、`/health` DB 非依存）は維持する。
- 分析原本: `docs/refactoring-instructions.md`。本書は実装上の判断が分かれる点の決定を記録する。

## 決定事項

### D1. DB サンプルは `/health/db`、handler は既存 `HealthHandler` を拡張する

CRUD サンプルではなく DB 死活確認（`SELECT 1` 相当の `PingContext`）を選ぶ。マイグレーションやサンプルテーブルを持ち込まず、最小コードで handler→service→repository→db の全連結と DI を実演できるため。handler は「health という1リソースに1ハンドラ」とし、新規 `DBHealthHandler` は作らず既存 `HealthHandler` に DB 用 service インターフェースと `GetDB` メソッドを追加する。

- repository へのインターフェース（`HealthPinger`）は利用側の service に定義する（Go の consumer-side interface イディオム）。
- Ping タイムアウトは業務判断として service に定数（2秒）で持つ。
- `db.OpenPostgres` がエラーを返した場合は `log.Fatalf` で終了してよい。DSN 不正は設定ミスであり「DB 未起動」とは別。DB 未起動は Ping 時に検出され 503 になるため、起動時 DB 非必須の要件は保たれる。`main.go` に `defer pool.Close()` を置く。

### D2. CORS は backend 側（Echo ミドルウェア）で対応する

代替案の Vite `server.proxy` は Go 側変更ゼロだが、エンドポイント追加のたびに proxy 設定が増える。テンプレートとしては backend 側 CORS が拡張に強い。許可オリジンは `CORS_ALLOW_ORIGINS`（カンマ区切り、デフォルト `http://localhost:5173`）として config に集約する。

### D3. ミドルウェア登録は router の責務とする

#4（composition root 移動）の「router はルート登録専任」は「Echo インスタンスの構成 + ルート登録」と定義する。router が config パッケージへ依存しないよう、許可オリジンは `[]string` 引数で受け取る。依存の生成（config → db → repository → service → handler）は `cmd/server/main.go` に置く。rules 文書の「router: 依存の配線」記述は「main: 依存の配線（composition root）、router: Echo 構成とルート登録」へ更新する。

### D4. go.mod は backend/ へ移動し、module 名を `…/backend` にする

module を `github.com/sudabon/lightweight_monorepo_go_project_template/backend` にすると、既存の import パス文字列（`…/backend/internal/…`）が一切変わらず、diff が go.mod の移動と module 行のみで済む。現状でも動作はするため優先度は中とし、実施順序は最後に置く（いつ実施してもよい独立作業）。

### D5. dotenv ライブラリは導入しない

`.env` 問題（README が `cp .env.example .env` を指示するが誰も読まない）は、ライブラリ追加ではなく README の手順修正で解消する（環境変数を明示して起動）。`.env.example` は「必要な環境変数のリファレンス」と位置づけを明記する。軽量さ優先。

### D6. frontend テストは Vitest + jsdom、フック抽出とセットで導入する

`App.tsx` から `useHealth`（`LoadState`・`refresh`・`AbortController` 含む）を `src/hooks/` へ抽出し、`vi.mock` で API をモックしたフックテストを1本置く。`vite.config.ts` に `test: { environment: "jsdom" }` を追加する（これがないと DOM 環境がなく失敗する）。Playwright は従来通り「必要になった時点で導入」。

### D7. 型ガードは `in` 演算子の絞り込みで書く

`response.json() as Promise<HealthResponse>` を廃し、`unknown` で受けて `isHealthResponse` で絞り込む。`in` 演算子による narrowing を使えば `as` を一切使わずに書ける（frontend 規約「unknown + 型ガード」の見本）。

### D8. 検証方法の注意

- doc comment の欠落は `go vet` では検出できない。受け入れはレビューで確認する（機械検査したい場合は `revive` の `exported` ルール）。
- `/health/db` の動作確認は `DATABASE_URL` を明示して起動する（D5 の通り `.env` は読まれない）。

## Risks / Trade-offs

- **module path 変更（D4）は BREAKING**: 既存フォークが本テンプレートの Go パッケージを import している場合に影響する。アプリケーションテンプレートのため外部 import はほぼ想定されないが、README のテンプレート利用手順で明示的に案内する。
- **CORS のデフォルト許可オリジン**: `http://localhost:5173` を許可するのは開発用デフォルトであり、本番では `CORS_ALLOW_ORIGINS` の明示が必要。環境変数表に用途を明記することで緩和する。
- **`/health/db` の追加は「DB なしでも起動する」性質を変えない**: OpenPostgres の `log.Fatalf` は DSN 不正（設定ミス）のみに限定し、DB 未起動は 503 で表現する（D1）。

## Migration Plan

1. tasks の順序（composition root → DB 縦切り + .env → CORS → 規約整合 → frontend → go.mod 移動）で段階的に実施する。各段階で `go test ./...` / `pnpm build` が通る状態を保ち、1関心事1コミットで分割する。
2. go.mod 移動は独立作業のため最後に置く。module 行の変更のみで import 文字列は不変（D4）。
3. テンプレート利用者向けには README の「テンプレート利用時」節を更新し、リネーム対象（`backend/go.mod` の module path）を案内する。ロールバックは各コミットの revert で可能。
