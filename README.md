# monorepo_go_project_template

クリーンアーキテクチャに沿ったモノレポ Web アプリ（Go バックエンド）を実装するための**ベース／テンプレート**用リポジトリです。依存の向きとソース配置をテストで固定し、エージェント向けの規約をリポジトリ内にまとめています。

## 現状のリポジトリに含まれるもの

| 内容 | 説明 |
|------|------|
| **アーキテクチャテスト** | `backend/tests/arch/` に `go test` ベースの検証があります。Go のソースルートは **`backend/internal/todo_app`**、モジュールパスは **`github.com/sudabon/todo_app`** を想定しています（フォーク時にリネームしてください）。 |
| **依存ルール** | `dependency_rule_test.go` … `domain` 層が `app` / `interfaces` / `infra` へ内向き import しないことを `go/parser` で検査します。 |
| **ディレクトリ構造** | `structure_test.go` … ソースルート直下に **`domain` → `app` → `interfaces` → `infra`** の4層フォルダのみがあることを検証します（内側から外側の順）。 |
| **エージェント規約** | `.cursor/rules/`、`.claude/rules/`、`.codex/rules/` にアーキテクチャ・バックエンド／フロント規約・テスト・Git などを配置しています。 |

依存の方向の原則（`.claude/rules/architecture.md` と同趣旨）は次のとおりです。

```
interfaces → app → domain
infra      → app → domain
```

- **domain** は他レイヤーに依存しない（Go標準ライブラリ中心）。
- **app** は **domain** のみ参照する。
- **infra** は **app** の抽象（リポジトリ/port等）を実装する。
- **interfaces** は **app** のユースケースを呼び出す。

> Goでは `interface` が予約語のため、Interface層のフォルダ／パッケージ名は **`interfaces`**（複数形）としています。

## まだ含まれていないもの（目標構成）

アプリケーション本体（Echo の `backend/cmd/`・`backend/internal/`、フロントエンド、`package.json` など）は**未配置**です。導入後に目指すスタック・ディレクトリ・コマンドの全体像は **`AGENTS.md`** に記載しています（Go 1.23 / Echo / GORM v2、React / Vite / pnpm など）。

アーキテクチャテストは、対象パッケージ（`backend/internal/todo_app`）が未配置のあいだは `t.Skip()` でスキップされます。`AGENTS.md` の構成に合わせて4層フォルダと実装を追加すると、依存ルール・構造の検査が有効になります。

```bash
go test ./...          # 素の状態では arch テストはスキップ（緑）
```

## ドキュメントの読み方

- **人間・ツール共通のプロジェクト概要・コマンド例** → [`AGENTS.md`](./AGENTS.md)
- **依存方向・禁止事項の短い要約** → [`.claude/rules/architecture.md`](./.claude/rules/architecture.md)

---

この README はリポジトリの実ファイル構成に基づいています。ソースやモジュール名を追加したら、`todo_app` 参照（ソースルート名・モジュールパス）とテストの期待値を自分のプロジェクト名に合わせて更新してください。
