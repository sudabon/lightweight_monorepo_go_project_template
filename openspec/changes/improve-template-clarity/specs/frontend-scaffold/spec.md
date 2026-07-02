## ADDED Requirements

### Requirement: データ取得ロジックはカスタムフックに置く

画面コンポーネントからデータ取得・状態管理ロジックを分離し、カスタムフックとして `src/hooks/` に配置しなければならない（SHALL）。参照実装として `useHealth` を提供し、ロード状態（loading / ready / error）と再取得関数を返す。コンポーネントはフックの戻り値を表示するだけにする。

#### Scenario: App がフック経由でデータを取得する

- **WHEN** `src/App.tsx` を確認する
- **THEN** `useState` / `useEffect` / `fetchHealth` の直接使用がなく、`useHealth` フックの戻り値を表示している

#### Scenario: フックが hooks ディレクトリにある

- **WHEN** `src/hooks/` を確認する
- **THEN** `useHealth.ts` が存在し、`LoadState` 型・再取得関数・AbortController によるクリーンアップを含む

### Requirement: Vitest によるユニットテスト基盤

frontend は Vitest によるユニットテストを実行できなければならない（SHALL）。テストファイルは対象と同階層に `*.test.ts(x)` として配置し、参照実装として `useHealth` のテストを提供する。

#### Scenario: pnpm test が成功する

- **WHEN** `cd frontend && pnpm test` を実行する
- **THEN** Vitest が実行され、すべてのテストが成功する

#### Scenario: ビルドも引き続き成功する

- **WHEN** `cd frontend && pnpm build` を実行する
- **THEN** 型チェックとビルドが成功する（フック抽出・型ガード・Vitest 設定の追加後も維持される）

#### Scenario: フックの成功・失敗が検証されている

- **WHEN** `src/hooks/useHealth.test.ts` を確認する
- **THEN** API クライアントをモックし、成功時に `ready`、失敗時に `error` になることを検証している

### Requirement: API レスポンスの型ガード検証

API クライアントはレスポンス JSON を `unknown` として受け取り、型ガード関数で絞り込んでから返さなければならない（SHALL）。型アサーション（`as`）によるレスポンス型の断定をしてはならない（MUST NOT）。

#### Scenario: 型ガードで検証してから返す

- **WHEN** `src/lib/api.ts` を確認する
- **THEN** `isHealthResponse` 型ガードが存在し、`as HealthResponse` / `as Promise<HealthResponse>` によるレスポンス型の断定が残っていない

#### Scenario: 不正な形状のレスポンスはエラーになる

- **WHEN** API が期待と異なる形状の JSON を返す
- **THEN** `fetchHealth` は例外を投げ、UI はエラー状態を表示する
