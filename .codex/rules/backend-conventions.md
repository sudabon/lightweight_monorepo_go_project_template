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

## Domain層

- entity は struct で定義し、不変条件はコンストラクタ関数（`NewXxx`）で検証する。
- ドメインサービスはステートレスにする。
- 外部ライブラリに依存しない（標準ライブラリ中心）。

## Application層

- usecase は struct とし、単一のパブリックメソッド（`Execute` 等）を持つ。
- リポジトリは Go の interface として `app/repository/` で定義する。
- 外部サービスの機能インターフェースは `app/port/` で定義する（通知・外部API等）。
- DTO は struct で定義し、入出力を明確にする。
- 共通の結果型やエラー型は `app/common/` に配置する。

## Infrastructure層

- GORM モデルは `infra/persistence/` に配置し、domain の entity とは分離する。
- infra モデル → domain entity の変換関数を用意する。
- app の repository/port インターフェースの具象実装をこの層に配置する。
- 設定値は `infra/config/` に集約し、環境変数から取得する。
- ロギングは `infra/logging/` に集約し、`log/slog` を使用する。

## Interface層（interfaces）

- Echo のハンドラ（router）は薄く保ち、ロジックは controller に委譲する。
- controller は usecase を呼び出し、結果を viewmodel に変換して返す。
- 依存（usecase・リポジトリ等）はコンストラクタ注入で渡す。
- リクエストのバインド／バリデーションは interfaces 層で行う。
