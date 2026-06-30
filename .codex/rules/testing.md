# テスト方針

## バックエンド

- `service` のユニットテストを優先する。
- `service` の業務判断は unit test を書く。
- `handler` は主要なHTTPステータスとレスポンスJSONをテストする。
- table-driven テストを基本とする。
- repository は実DBが必要な場合のみ integration test として追加し、テスト用DB（testcontainers 等）または build tag で分離して検証する。
- `/health` はDB非依存を維持する。
- テストファイルは対象と同階層に `xxx_test.go` として配置する。
- テスト関数名は `TestXxx_条件_期待結果` の形式にする。
- 共通のセットアップは `TestMain` やヘルパー関数（`t.Helper()`）に集約する。

## フロントエンド

- ユニットテスト（Vitest）: カスタムフック、ユーティリティ関数、複雑なロジックを持つコンポーネントを対象にする。
- E2Eテスト（Playwright）: ユーザーフロー単位で `frontend/e2e/` に配置する。
- テストファイルは対象ファイルと同階層に `*.test.ts(x)` として配置する。
