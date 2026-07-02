---
globs: ["backend/**"]
---
# バックエンド コーディング規約

## 全般

- 公開（エクスポート）する型・関数には doc comment を付与する（`// Xxx は…` 形式）。
- 命名: エクスポートは PascalCase、非公開は camelCase、パッケージ名は短い小文字（アンダースコアなし）。
- `any`（`interface{}`）の濫用を避ける。必要な箇所のみに限定する。
- エラーは値として返し、`if err != nil` で必ずハンドリングする。握りつぶさない。
- キャンセル・タイムアウトのため `context.Context` を第1引数で伝播する。
- マジックナンバー・文字列は定数（`const`）として定義する。
- `gofmt` で整形し、`go vet` を通す。

## handler

- Echo の `Context` は handler 内に閉じ込める。
- service を呼び、HTTP ステータスと JSON レスポンスを返す。
- SQL、DBドライバ、環境変数アクセスを書かない。

## service

- 業務判断とユースケースの流れを置く。
- 永続化が必要な場合は repository を呼ぶ。
- HTTP 固有型に依存しない。
- 小規模なレスポンス型は `service` に置いてよい。ただし、HTTP専用の表現や画面都合の型が増えた場合は `handler` 側または `dto` に分離する。
- table-driven test を優先して追加する。

## repository

- DB・SQL・ドライバ固有処理を置く。
- `database/sql` の `*sql.DB` はコンストラクタで受け取る。
- handler から直接呼ばせない。

## config / db / router

- 環境変数は `config` に集約する。
- DB接続生成は `db` に置く。
- `db.OpenPostgres` は接続プールを初期化するだけで、接続確認は行わない。起動時にDB必須とする場合は、呼び出し側で `PingContext` を実行する。
- Echo の構成とルート登録は `router` に置く。
- 依存の配線（composition root）は `cmd/server/main.go` に置く。
