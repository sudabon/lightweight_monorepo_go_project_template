# テスト方針

## バックエンド

- ユニットテスト: domain, app 層を中心にテストする。リポジトリ/port は interface のモック実装で差し替える。
- table-driven テストを基本とする。
- 統合テスト: infra 層のDB操作等を、テスト用DB（testcontainers 等）または build tag で分離して検証する。
- テストファイルは対象と同階層に `xxx_test.go` として配置する。アーキテクチャテストは `backend/tests/arch/` に置く。
- テスト関数名は `TestXxx_条件_期待結果` の形式にする。
- 共通のセットアップは `TestMain` やヘルパー関数（`t.Helper()`）に集約する。

## フロントエンド

- ユニットテスト（Vitest）: カスタムフック、ユーティリティ関数、複雑なロジックを持つコンポーネントを対象にする。
- E2Eテスト（Playwright）: ユーザーフロー単位で `frontend/e2e/` に配置する。
- テストファイルは対象ファイルと同階層に `*.test.ts(x)` として配置する。
