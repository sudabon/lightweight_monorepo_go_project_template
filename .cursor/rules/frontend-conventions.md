---
globs: ["frontend/**"]
---
# フロントエンド コーディング規約

- `any` 型の使用は禁止。`unknown` を使い、型ガードで絞り込む。
- コンポーネントは関数コンポーネントのみ使用する。
- コンポーネント名は PascalCase、関数・変数は camelCase。
- API クライアントは `frontend/src/lib/` に集約する。
- API レスポンス型は `frontend/src/types/` に配置する。
- コンポーネントから直接 `fetch` を増やさない。
- TanStack Query、TailwindCSS、UIライブラリは必要になった時点で導入する。
- マジックナンバー・文字列リテラルは定数として定義する。
- フロントエンドのコマンドは `frontend/` ディレクトリ内で実行する。
