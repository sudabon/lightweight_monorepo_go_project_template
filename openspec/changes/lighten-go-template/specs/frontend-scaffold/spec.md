## ADDED Requirements

### Requirement: React + TypeScript + Vite の最小スキャフォルド

frontend は React + TypeScript + Vite を前提とした最小スキャフォルドを提供しなければならない（SHALL）。`package.json`、`vite.config.ts`、`index.html`、`src/main.tsx`、`src/App.tsx` と、`components/`、`pages/`、`hooks/`、`lib/`、`types/` のディレクトリを持つ。

#### Scenario: frontend スキャフォルドが存在する

- **WHEN** `frontend/` を確認する
- **THEN** `package.json`、`vite.config.ts`、`index.html`、`src/main.tsx`、`src/App.tsx` と `src/{components,pages,hooks,lib,types}/` が存在する

#### Scenario: ビルドが成功する

- **WHEN** `cd frontend && pnpm install && pnpm build` を実行する
- **THEN** TypeScript の型エラーなくプロダクションビルドが成功する

### Requirement: API呼び出しと型定義の集約

API呼び出しは `src/lib`（例: `src/lib/api.ts`）に集約しなければならず、画面コンポーネントから直接 `fetch` を乱用してはならない（MUST NOT）。APIレスポンス型は `src/types/` に配置する。TanStack Query・Shadcn/ui は任意導入とし、軽量版では必須にしない。

#### Scenario: APIクライアントが lib に集約される

- **WHEN** frontend のAPI呼び出しコードを確認する
- **THEN** API呼び出しは `src/lib/api.ts` 経由で行われ、レスポンス型は `src/types/` に定義されている
