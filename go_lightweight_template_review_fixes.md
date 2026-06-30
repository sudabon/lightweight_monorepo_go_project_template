# lightweight_monorepo_go_project_template 修正指示書

対象リポジトリ:

- https://github.com/sudabon/lightweight_monorepo_go_project_template

## 目的

Go + React の軽量モノレポテンプレートとして、現状の `handler / service / repository` という3責務構成は維持する。

厳密な Clean Architecture には戻さず、以下を目的に最小限の修正を行う。

- テンプレートとしての信頼性を上げる
- エージェント開発時に責務が崩れにくくする
- 新規プロジェクトへ転用したときの初期品質を上げる
- README / AGENTS.md / 実装のズレを減らす

## 修正方針

このリポジトリは軽量版テンプレートであるため、過度な抽象化は避ける。

特に以下は行わない。

- `domain / app / interfaces / infra` の4層構成へ戻さない
- repository interface を先回りして大量に定義しない
- DIコンテナやwireなどを最初から導入しない
- TanStack Query / TailwindCSS / UIライブラリを必須化しない

## 優先度A: 最低限のテストを追加する

### 背景

README では `go test ./...` を実行コマンドとして案内しているが、現状では最低限のテストが不足している。

テンプレート利用者やエージェントが変更したときに、最初から確認できるテストが必要。

### 追加するファイル

```text
backend/internal/service/health_service_test.go
backend/internal/handler/health_handler_test.go
```

### 実装内容

#### 1. `HealthService.Status()` のテスト

`service.NewHealthService().Status()` が以下を返すことを確認する。

```json
{"status":"ok"}
```

確認観点:

- `Status` が `"ok"` であること
- panicしないこと
- DB接続を必要としないこと

#### 2. `HealthHandler.Get()` のテスト

Echo のテスト機能を使い、`GET /health` 相当のハンドラ呼び出しをテストする。

確認観点:

- HTTPステータスが `200 OK` であること
- レスポンスJSONが `{"status":"ok"}` であること
- handler からDB接続が不要であること

### 完了条件

```bash
cd backend
go test ./...
go vet ./...
gofmt -l .
```

が成功すること。

`gofmt -l .` は未整形ファイルが出ないこと。

## 優先度B: DTO / response 型の置き場所を検討する

### 背景

現状の `handler` は `service.HealthStatus` に依存している。

軽量テンプレートとしては許容範囲だが、今後エンドポイントが増えた場合、HTTPレスポンス型と service の戻り値が密結合になりやすい。

### 推奨修正

以下のいずれかを選ぶ。

#### 案1: 現状維持

小さなテンプレートとしては最も軽い。

```text
handler -> service
```

の依存を許容する。

この場合、READMEまたはAGENTS.mdに以下の方針を追記する。

```text
小規模なレスポンス型は service に置いてよい。
ただし、HTTP専用の表現や画面都合の型が増えた場合は handler 側または dto に分離する。
```

#### 案2: `internal/dto` を追加する

以下を追加する。

```text
backend/internal/dto/
  health.go
```

`HealthStatus` を `dto.HealthStatus` として配置する。

依存関係は以下とする。

```text
handler -> dto
service -> dto
```

### 推奨判断

軽量テンプレートとしては **案1でよい**。

ただし、MR向け業務システムなど実案件に転用する場合は、早めに `internal/dto` または `internal/model` を追加してよい。

## 優先度B: README / AGENTS.md にテスト方針を追記する

### 背景

軽量テンプレートでは、構造をテストで縛りすぎない代わりに、最低限のテスト方針を明記しておくとエージェントが壊しにくい。

### 追記内容

README または AGENTS.md に以下を追加する。

```text
## テスト方針

- service の業務判断は unit test を書く。
- handler は主要なHTTPステータスとレスポンスJSONをテストする。
- repository は実DBが必要な場合のみ integration test として追加する。
- `/health` はDB非依存を維持する。
```

## 優先度C: DB接続の扱いを明記する

### 背景

`backend/internal/db/postgres.go` は `sql.Open("pgx", databaseURL)` を返しているが、`sql.Open` は即時接続確認を行わない。

テンプレート利用者が「DB接続成功」と誤解しないように、必要に応じて `PingContext` の扱いをREADMEに明記するとよい。

### 推奨追記

```text
`db.OpenPostgres` は接続プールを初期化するだけで、接続確認は行わない。
起動時にDB必須とする場合は、呼び出し側で `PingContext` を実行する。
ただし、初期テンプレートでは `/health` をDB非依存にするため、起動時DB接続確認は行わない。
```

## 優先度C: frontend API クライアントに末尾スラッシュ除去を追加する

### 背景

Python版では `VITE_API_BASE_URL` の末尾スラッシュを除去しているが、Go版では除去していない。

Go版でも揃えると、以下のような設定ミスに強くなる。

```env
VITE_API_BASE_URL=http://localhost:8080/
```

### 修正案

`frontend/src/lib/api.ts` を以下の方針にする。

```ts
const API_BASE_URL = (import.meta.env.VITE_API_BASE_URL ?? "http://localhost:8080").replace(/\/$/, "");
```

## 優先度C: frontend API レスポンスの型ガードを追加するか検討する

### 背景

Python版では API レスポンスを `unknown` として受け取り、型ガードで検証している。

Go版は `response.json() as Promise<HealthResponse>` で型アサーションしている。

テンプレートとして安全側に寄せるなら、Go版 frontend にも型ガードを追加する。

### 推奨判断

軽量さ重視なら現状維持でよい。

ただし、Python版と方針を揃えるなら、以下を追加する。

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

## 優先度C: docker-compose のサービス名を揃えるか検討する

### 背景

Go版ではDBサービス名が `db`、Python版では `postgres` になっている。

単体では問題ないが、両テンプレートを横並びで管理する場合、名前が揃っている方が理解しやすい。

### 推奨判断

Go版のREADMEでは `docker compose up -d db` と案内しているため、現状維持でよい。

両方を統一するなら、どちらかに寄せる。

推奨は以下。

```text
services:
  db:
```

Python版も `db` に寄せると、READMEの説明が揃う。

## 変更後に実行する確認コマンド

```bash
cd backend
go mod tidy
go test ./...
go vet ./...
gofmt -l .
```

```bash
cd frontend
pnpm install
pnpm build
```

## 最終チェックリスト

- [ ] `/health` はDB非依存で動作する
- [ ] `go test ./...` が成功する
- [ ] `go vet ./...` が成功する
- [ ] `gofmt -l .` で未整形ファイルが出ない
- [ ] `handler` からDBを直接触っていない
- [ ] `handler` にSQLを書いていない
- [ ] 環境変数アクセスは `backend/internal/config` に集約されている
- [ ] frontend コンポーネントから直接 `fetch` を増やしていない
- [ ] README と AGENTS.md の内容が実装と矛盾していない
