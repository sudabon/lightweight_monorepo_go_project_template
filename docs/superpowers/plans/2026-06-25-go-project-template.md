# Go版プロジェクトテンプレート Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** `monorepo_python_project_template` のGo言語版スケルトンを本リポジトリに作成する（アーキテクチャテスト + エージェント規約 + ドキュメント。アプリ本体は未配置）。

**Architecture:** Clean Architectureの4層（domain / app / interfaces / infra）を、Go標準の `go/parser` を用いたアーキテクチャテストで固定する。技術スタック（Echo / GORM v2 / golang-migrate）は規約ドキュメントに「採用予定」として記載し、実コードは作らない。対象パッケージ未配置の素の状態では arch テストは `t.Skip()` で緑になる。

**Tech Stack:** Go 1.23、標準ライブラリ（`go/parser`, `go/token`, `os`, `path/filepath`）、`go test`。

## Global Constraints

- Go バージョン: `go 1.23`（go.mod の `go` ディレクティブ）。
- モジュールパス: `github.com/sudabon/todo_app`（フォーク時にリネームする旨を注記）。
- 層命名: `domain` / `app` / `interfaces`（複数形・Go予約語回避） / `infra`。
- 依存方向: `interfaces → app → domain`、`infra → app → domain`。
- アーキテクチャテストは外部依存ゼロ（標準ライブラリのみ）。
- 素のテンプレートで `go test ./...` / `go vet ./...` が緑、`gofmt -l backend` が無出力であること。
- 規約ドキュメントは日本語。コミットメッセージは日本語 + プレフィックス（`feat:` `fix:` `refactor:` `test:` `docs:` `chore:`）。
- 元テンプレート（コピー元）の絶対パス: `/Users/suda/workspace/sudabon/monorepo_python_project_template`。

---

## ファイル構成

| ファイル | 責務 |
|----------|------|
| `go.mod` | モジュール定義。archテストのコンパイル/実行に必要。 |
| `backend/tests/arch/structure_test.go` | ソースルート直下が4層のみであることを検査。 |
| `backend/tests/arch/dependency_rule_test.go` | domain層が外向きimportを持たないことを検査。 |
| `AGENTS.md` | プロジェクト概要・スタック・ディレクトリ・コマンド。 |
| `CLAUDE.md` | `@AGENTS.md` の1行。 |
| `README.md` | テンプレートの説明。 |
| `.claude/rules/*.md` | エージェント規約7ファイル。 |
| `.cursor/rules/*.md` / `.codex/rules/*.md` | 上記のうち6ファイルのミラー。 |
| `.claude/skills/**` | 言語非依存スキル2件を複製。 |

---

### Task 1: Goモジュール + ソース構造テスト

**Files:**
- Create: `go.mod`
- Create: `backend/tests/arch/structure_test.go`

**Interfaces:**
- Produces: パッケージ `arch`（`backend/tests/arch`）。定数 `sourceRoot string`（このパッケージ内で共有、Task 2が再利用）、`var expectedLayers []string`。

- [ ] **Step 1: go.mod を作成する**

`go.mod`:

```
module github.com/sudabon/todo_app

go 1.23
```

- [ ] **Step 2: 構造テストを書く**

`backend/tests/arch/structure_test.go`:

```go
package arch

import (
	"os"
	"sort"
	"testing"
)

// sourceRoot は Clean Architecture のソースルート。このテストパッケージ
// （backend/tests/arch）からの相対パス。フォーク時にリネームする。
const sourceRoot = "../../internal/todo_app"

// expectedLayers は sourceRoot 直下に許可される唯一のディレクトリ群
// （内側から外側の順）。
var expectedLayers = []string{"domain", "app", "interfaces", "infra"}

// TestSourceStructure は sourceRoot 直下が4層のみで構成されることを検査する。
// 対象パッケージが未配置の場合はスキップする（テンプレート素の状態）。
func TestSourceStructure(t *testing.T) {
	if _, err := os.Stat(sourceRoot); os.IsNotExist(err) {
		t.Skipf("source root %q not found; add the application package to enable this check", sourceRoot)
	}

	entries, err := os.ReadDir(sourceRoot)
	if err != nil {
		t.Fatalf("read source root: %v", err)
	}

	found := map[string]bool{}
	for _, e := range entries {
		if e.IsDir() {
			found[e.Name()] = true
		}
	}

	for _, layer := range expectedLayers {
		if !found[layer] {
			t.Errorf("missing %q layer directory under %s", layer, sourceRoot)
		}
	}

	expected := map[string]bool{}
	for _, layer := range expectedLayers {
		expected[layer] = true
	}
	var unexpected []string
	for name := range found {
		if !expected[name] {
			unexpected = append(unexpected, name)
		}
	}
	sort.Strings(unexpected)
	if len(unexpected) > 0 {
		t.Errorf("source root should only contain Clean Architecture layers; unexpected: %v", unexpected)
	}
}
```

- [ ] **Step 3: スキップを確認する**

Run: `go test ./backend/tests/arch/ -run TestSourceStructure -v`
Expected: `--- SKIP: TestSourceStructure` （source root not found）

- [ ] **Step 4: テストが違反を検出することを確認する（フィクスチャ）**

```bash
mkdir -p backend/internal/todo_app/domain backend/internal/todo_app/app backend/internal/todo_app/interfaces backend/internal/todo_app/infra
go test ./backend/tests/arch/ -run TestSourceStructure -v
```
Expected: PASS（4層が揃い、想定外なし）

```bash
mkdir -p backend/internal/todo_app/extra
go test ./backend/tests/arch/ -run TestSourceStructure -v
```
Expected: FAIL（`unexpected: [extra]`）

- [ ] **Step 5: フィクスチャを削除しスキップに戻ることを確認する**

```bash
rm -rf backend/internal
go test ./backend/tests/arch/ -run TestSourceStructure -v
```
Expected: `--- SKIP: TestSourceStructure`

- [ ] **Step 6: フォーマット・vetを確認する**

Run: `gofmt -l backend && go vet ./...`
Expected: 無出力 / exit 0

- [ ] **Step 7: コミット**

```bash
git add go.mod backend/tests/arch/structure_test.go
git commit -m "test: Goモジュールとソース構造のアーキテクチャテストを追加"
```

---

### Task 2: 依存ルールテスト

**Files:**
- Create: `backend/tests/arch/dependency_rule_test.go`

**Interfaces:**
- Consumes: `sourceRoot`（Task 1で定義、同一パッケージ `arch`）。
- Produces: `TestDomainLayerDependencies`、ヘルパー `layerAfterSourceRoot(importPath string) (string, bool)`、定数 `sourceRootName`。

- [ ] **Step 1: 依存ルールテストを書く**

`backend/tests/arch/dependency_rule_test.go`:

```go
package arch

import (
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// sourceRootName はソースルートのディレクトリ/パッケージ名。依存検査は
// import パス中のこのセグメント直後の層名をキーに行う。フォーク時にリネーム。
const sourceRootName = "todo_app"

// forbiddenForDomain は domain 層が import してはならない層。
var forbiddenForDomain = map[string]bool{"app": true, "interfaces": true, "infra": true}

// TestDomainLayerDependencies は domain 層が外側の層へ内向き import しない
// ことを検査する。domain ディレクトリが無い場合はスキップする。
func TestDomainLayerDependencies(t *testing.T) {
	domainPath := filepath.Join(sourceRoot, "domain")
	if _, err := os.Stat(domainPath); os.IsNotExist(err) {
		t.Skipf("domain layer %q not found; add the application package to enable this check", domainPath)
	}

	fset := token.NewFileSet()
	var violations []string

	err := filepath.WalkDir(domainPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || !strings.HasSuffix(path, ".go") {
			return nil
		}
		file, perr := parser.ParseFile(fset, path, nil, parser.ImportsOnly)
		if perr != nil {
			return perr
		}
		for _, imp := range file.Imports {
			layer, ok := layerAfterSourceRoot(strings.Trim(imp.Path.Value, `"`))
			if ok && forbiddenForDomain[layer] {
				rel, _ := filepath.Rel(domainPath, path)
				violations = append(violations,
					rel+": domain layer cannot import from "+layer+" layer")
			}
		}
		return nil
	})
	if err != nil {
		t.Fatalf("walk domain layer: %v", err)
	}

	if len(violations) > 0 {
		t.Errorf("dependency rule violations:\n%s", strings.Join(violations, "\n"))
	}
}

// layerAfterSourceRoot は import パス中で sourceRootName の直後に来る層名を返す。
// モジュールパス自体に sourceRootName と同名のセグメントが含まれる場合
// （例: module github.com/sudabon/todo_app 配下の todo_app パッケージ）に
// 誤検出しないよう、最後に出現する sourceRootName を基準にする。
func layerAfterSourceRoot(importPath string) (string, bool) {
	segs := strings.Split(importPath, "/")
	for i := len(segs) - 1; i >= 0; i-- {
		if segs[i] == sourceRootName && i+1 < len(segs) {
			return segs[i+1], true
		}
	}
	return "", false
}
```

- [ ] **Step 2: スキップを確認する**

Run: `go test ./backend/tests/arch/ -run TestDomainLayerDependencies -v`
Expected: `--- SKIP: TestDomainLayerDependencies`

- [ ] **Step 3: 違反検出を確認する（フィクスチャ）**

```bash
mkdir -p backend/internal/todo_app/domain
cat > backend/internal/todo_app/domain/bad.go <<'EOF'
package domain

import _ "github.com/sudabon/todo_app/backend/internal/todo_app/app/usecase"
EOF
go test ./backend/tests/arch/ -run TestDomainLayerDependencies -v
```
Expected: FAIL（`bad.go: domain layer cannot import from app layer`）

- [ ] **Step 4: フィクスチャを削除しスキップに戻ることを確認する**

```bash
rm -rf backend/internal
go test ./backend/tests/arch/ -run TestDomainLayerDependencies -v
```
Expected: `--- SKIP: TestDomainLayerDependencies`

- [ ] **Step 5: 全体検証**

Run: `gofmt -l backend && go vet ./... && go test ./...`
Expected: 無出力 / 全テスト ok（archはskip）

- [ ] **Step 6: コミット**

```bash
git add backend/tests/arch/dependency_rule_test.go
git commit -m "test: domain層の依存ルールのアーキテクチャテストを追加"
```

---

### Task 3: AGENTS.md と CLAUDE.md

**Files:**
- Create: `AGENTS.md`
- Create: `CLAUDE.md`

- [ ] **Step 1: AGENTS.md を作成する**

`AGENTS.md`:

````markdown
# AGENTS.md

このファイルはClaude Codeがリポジトリを操作する際のガイドラインである。
詳細なルールは `.claude/rules/` に分割して配置している。

## プロジェクト概要

モノレポ構成のWebアプリケーション。

- **バックエンド** (`backend/`): Go 1.23 / Echo。Clean Architectureに基づく4層構成。
- **フロントエンド** (`frontend/`): React / TypeScript / Vite / TailwindCSS。

## ディレクトリ構成

```
├── backend/                       # バックエンド
│   ├── cmd/
│   │   └── server/
│   │       └── main.go            # サーバプロセス管理（エントリポイント）
│   ├── internal/
│   │   └── todo_app/              # ソースルート（フォーク時にリネーム）
│   │       ├── app/               # Application層
│   │       │   ├── common/        # 共通の結果型・エラー型
│   │       │   ├── dto/           # Data Transfer Object（層間のデータ受け渡し）
│   │       │   ├── repository/    # データ永続化のインターフェース定義
│   │       │   ├── port/          # 外部サービス機能のインターフェース定義
│   │       │   └── usecase/       # ビジネスロジックの実行フロー
│   │       ├── domain/            # Domain層
│   │       │   ├── entity/        # エンティティ（ドメインモデル）
│   │       │   └── service/       # ドメインサービス（エンティティ横断ロジック）
│   │       ├── infra/             # Infrastructure層
│   │       │   ├── cli/           # CLIコマンド
│   │       │   ├── config/        # 設定（環境変数）
│   │       │   ├── notification/  # 通知（メール・Slack等）
│   │       │   ├── logging/       # ロギング（slog）
│   │       │   ├── persistence/   # データ永続化（GORMモデル・リポジトリ実装）
│   │       │   └── web/           # Web関連（Echoルーター・ミドルウェア）
│   │       └── interfaces/        # Interface層
│   │           ├── controller/    # 入力リクエスト処理
│   │           ├── presenter/     # 出力変換（domain固有型をviewmodel固有型へ）
│   │           └── viewmodel/     # レスポンス整形用ViewModel
│   ├── migrations/                # golang-migrate のSQLファイル
│   └── tests/
│       └── arch/                  # クリーンアーキテクチャのテスト
├── frontend/                      # フロントエンド
│   ├── src/
│   │   ├── components/            # UIコンポーネント（Shadcn/ui ベース）
│   │   ├── pages/                 # ページコンポーネント
│   │   ├── hooks/                 # カスタムフック
│   │   ├── lib/                   # ユーティリティ・API クライアント
│   │   ├── types/                 # 型定義
│   │   └── routes/                # React Router ルート定義
│   ├── e2e/                       # Playwright E2Eテスト
│   ├── package.json
│   └── vite.config.ts
├── go.mod
```

## 技術スタック

### バックエンド

- **言語**: Go 1.23+
- **フレームワーク**: Echo (v4)
- **パッケージ管理**: go modules
- **ORM**: GORM v2
- **DB**: PostgreSQL
- **マイグレーション**: golang-migrate
- **テスト**: go test（table-driven）
- **リンター**: golangci-lint / gofmt
- **静的解析**: go vet
- **DI**: コンストラクタ注入（手動配線。必要に応じ google/wire）

### フロントエンド

- **言語**: TypeScript
- **フレームワーク**: React
- **ビルドツール**: Vite
- **パッケージ管理**: pnpm
- **UIライブラリ**: Shadcn/ui + TailwindCSS
- **データ取得**: TanStack Query
- **ルーティング**: React Router
- **テスト**: Vitest（ユニット）/ Playwright（E2E）
- **リンター**: ESLint

## コマンド

### バックエンド

```bash
# 開発環境
go mod download                # 依存パッケージのダウンロード
go mod tidy                    # 依存の整理
go run ./backend/cmd/server    # 開発サーバー起動

# テスト
go test ./...                  # テスト全体実行
go test ./... -run TestName    # 特定テストの実行
go test ./... -failfast        # 最初の失敗で停止
go test ./backend/tests/arch/  # アーキテクチャテストのみ

# リント・静的解析・フォーマット
golangci-lint run              # リントチェック
go vet ./...                   # 静的解析
gofmt -l .                     # 未整形ファイルの一覧
gofmt -w .                     # フォーマット適用

# マイグレーション（golang-migrate）
migrate -path backend/migrations -database "$DATABASE_URL" up    # 最新まで適用
migrate -path backend/migrations -database "$DATABASE_URL" down 1 # 1つ前にロールバック
migrate create -ext sql -dir backend/migrations -seq <name>      # マイグレーション作成
```

### フロントエンド

```bash
cd frontend
pnpm install                   # 依存パッケージのインストール
pnpm dev                       # 開発サーバー起動
pnpm test                      # Vitest ユニットテスト実行
pnpm test:e2e                  # Playwright E2Eテスト実行
pnpm lint                      # ESLint チェック
pnpm build                     # プロダクションビルド
```

## 注意事項

- DBのURL等はすべて環境変数から取得する（`infra/config/` に集約）。
- バックエンド・フロントエンド間の型定義の乖離に注意する。APIスキーマを信頼の源泉とする。
- `interface` はGoの予約語のため、Interface層のパッケージ名は `interfaces`（複数形）とする。
````

- [ ] **Step 2: CLAUDE.md を作成する**

`CLAUDE.md`:

```
@AGENTS.md
```

- [ ] **Step 3: 検証**

Run: `test -f AGENTS.md && test -f CLAUDE.md && head -1 CLAUDE.md`
Expected: `@AGENTS.md`

- [ ] **Step 4: コミット**

```bash
git add AGENTS.md CLAUDE.md
git commit -m "docs: Go版のAGENTS.md・CLAUDE.mdを追加"
```

---

### Task 4: README.md

**Files:**
- Modify: `README.md`（既存のGo用READMEを上書き）

- [ ] **Step 1: README.md を上書きする**

`README.md`:

````markdown
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
````

- [ ] **Step 2: 検証**

Run: `grep -q "monorepo_go_project_template" README.md && grep -q "interfaces → app → domain" README.md && echo OK`
Expected: `OK`

- [ ] **Step 3: コミット**

```bash
git add README.md
git commit -m "docs: READMEをGo版テンプレートの内容に更新"
```

---

### Task 5: .claude/rules/ 規約7ファイル

**Files:**
- Create: `.claude/rules/architecture.md`
- Create: `.claude/rules/backend-conventions.md`
- Create: `.claude/rules/frontend-conventions.md`
- Create: `.claude/rules/testing.md`
- Create: `.claude/rules/git.md`
- Create: `.claude/rules/templates.md`
- Create: `.claude/rules/codex-guidelines.md`

- [ ] **Step 1: architecture.md**

`.claude/rules/architecture.md`:

```markdown
# アーキテクチャ原則

## 依存の方向（厳守）

​```
interfaces → app → domain
infra      → app → domain
​```

- **domain** は他のどの層にも依存しない。Go標準ライブラリのみ使用可。
- **app** は domain のみに依存する。infra や interfaces を import しない。
- **infra** は app のリポジトリ/port インターフェースを実装する。
- **interfaces** は app の usecase を呼び出す。domain を直接操作しない。

## 禁止事項

- domain 層で GORM・Echo・その他外部ライブラリを import しない
- usecase 内で具象リポジトリや具象port を直接生成しない（コンストラクタでDIする）
- controller から entity を直接返さない（必ず viewmodel に変換する）
- infra の具象型を interfaces から直接 import しない
- app 層の port インターフェースに infra の実装詳細を漏らさない
```

（注: 上記コードフェンス内の依存方向ブロックは、実ファイルでは通常の三連バッククォートで囲むこと。ゼロ幅文字 `​` は使わない。）

- [ ] **Step 2: backend-conventions.md**

`.claude/rules/backend-conventions.md`:

```markdown
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
```

- [ ] **Step 3: frontend-conventions.md（Python版から変更なし）**

`.claude/rules/frontend-conventions.md`:

```markdown
---
globs: ["frontend/**"]
---
# フロントエンド コーディング規約

- `any` 型の使用は禁止。`unknown` を使い、型ガードで絞り込む。
- コンポーネントは関数コンポーネントのみ使用する。
- コンポーネント名は PascalCase、関数・変数は camelCase。
- UIコンポーネントは Shadcn/ui をベースに使用する。独自スタイルは TailwindCSS で適用する。
- サーバー状態の取得・キャッシュは TanStack Query を使用する。`useEffect` + `fetch` で代用しない。
- API クライアントは `frontend/src/lib/` に集約する。
- マジックナンバー・文字列リテラルは定数として定義する。
- フロントエンドのコマンドは `frontend/` ディレクトリ内で実行する。
```

- [ ] **Step 4: testing.md**

`.claude/rules/testing.md`:

```markdown
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
```

- [ ] **Step 5: git.md（Python版から変更なし）**

`.claude/rules/git.md`:

```markdown
# Git ルール

- コミットメッセージは日本語で記述する。
- プレフィックス: `feat:`, `fix:`, `refactor:`, `test:`, `docs:`, `chore:`
- 1コミット1関心事。大きな変更は分割する。
- `.env` ファイルはコミットしない。`.env.example` を参考にする。
```

- [ ] **Step 6: templates.md**

`.claude/rules/templates.md`:

````markdown
# ファイル作成時のテンプレート

## 新規 entity

```go
package entity

import "errors"

// UserID はユーザーの識別子。
type UserID string

// User はユーザーエンティティ。
type User struct {
	ID   UserID
	Name string
}

// NewUser は不変条件を検証して User を生成する。
func NewUser(id UserID, name string) (*User, error) {
	if name == "" {
		return nil, errors.New("name must not be empty")
	}
	return &User{ID: id, Name: name}, nil
}
```

## 新規 usecase

```go
package usecase

import "context"

// XxxUseCase は Xxx のユースケース。
type XxxUseCase struct {
	xxxRepo repository.XxxRepository
}

// NewXxxUseCase は XxxUseCase を生成する。
func NewXxxUseCase(xxxRepo repository.XxxRepository) *XxxUseCase {
	return &XxxUseCase{xxxRepo: xxxRepo}
}

// Execute はユースケースを実行する。
func (u *XxxUseCase) Execute(ctx context.Context, in XxxInput) (XxxOutput, error) {
	// ...
}
```

## 新規 repository インターフェース

```go
package repository

import "context"

// XxxRepository は Xxx の永続化インターフェース。
type XxxRepository interface {
	FindByID(ctx context.Context, id entity.XxxID) (*entity.Xxx, error)
	Save(ctx context.Context, x *entity.Xxx) error
}
```

## 新規 router / handler（Echo）

```go
package web

import (
	"github.com/labstack/echo/v4"
)

// RegisterXxxRoutes は Xxx 関連のルートを登録する。
func RegisterXxxRoutes(e *echo.Echo, c *controller.XxxController) {
	g := e.Group("/xxx")
	g.GET("", c.List)
}
```

## 新規 port インターフェース

```go
package port

import "context"

// XxxPort は外部サービス Xxx のインターフェース。
type XxxPort interface {
	Send(ctx context.Context, msg XxxMessage) error
}
```
````

- [ ] **Step 7: codex-guidelines.md（Python版から変更なし、.claude のみ）**

`.claude/rules/codex-guidelines.md`:

```markdown
## いつ使うか

### 必須（自動実行）
- 資料・成果物を作成した後 → Codexにレビューを依頼
- 重要な意思決定の前 → Codexの見解を聞く

### 推奨
- アプローチが2回以上失敗したとき
- 複数の選択肢で迷ったとき
- 専門外の領域で自信がないとき
```

- [ ] **Step 8: 検証**

Run: `ls .claude/rules/ | sort | tr '\n' ' '`
Expected: `architecture.md backend-conventions.md codex-guidelines.md frontend-conventions.md git.md templates.md testing.md`

Run: `grep -q "interfaces → app → domain" .claude/rules/architecture.md && grep -q "GORM" .claude/rules/backend-conventions.md && echo OK`
Expected: `OK`

- [ ] **Step 9: コミット**

```bash
git add .claude/rules/
git commit -m "docs: Go版のエージェント規約（.claude/rules）を追加"
```

---

### Task 6: .cursor/rules と .codex/rules へミラー

`.claude/rules/` のうち `codex-guidelines.md` を除く6ファイルを、`.cursor/rules/` と `.codex/rules/` に複製する（Python版と同じ構成）。

**Files:**
- Create: `.cursor/rules/{architecture,backend-conventions,frontend-conventions,testing,git,templates}.md`
- Create: `.codex/rules/{architecture,backend-conventions,frontend-conventions,testing,git,templates}.md`

- [ ] **Step 1: ミラーする**

```bash
mkdir -p .cursor/rules .codex/rules
for f in architecture backend-conventions frontend-conventions testing git templates; do
  cp ".claude/rules/$f.md" ".cursor/rules/$f.md"
  cp ".claude/rules/$f.md" ".codex/rules/$f.md"
done
```

- [ ] **Step 2: 検証**

Run: `diff <(ls .cursor/rules) <(ls .codex/rules) && ls .cursor/rules | wc -l`
Expected: 差分なし / `6`

Run: `diff .claude/rules/architecture.md .cursor/rules/architecture.md && echo SAME`
Expected: `SAME`

- [ ] **Step 3: コミット**

```bash
git add .cursor/rules/ .codex/rules/
git commit -m "docs: 規約を.cursor/.codexにミラー"
```

---

### Task 7: スキルの複製

言語非依存のワークフロースキル2件を元テンプレートから複製する。

**Files:**
- Create: `.claude/skills/opsx-codex-apply/**`
- Create: `.claude/skills/openspec-adr-generator/**`

- [ ] **Step 1: 複製する**

```bash
mkdir -p .claude/skills
SRC=/Users/suda/workspace/sudabon/monorepo_python_project_template/.claude/skills
cp -R "$SRC/opsx-codex-apply" .claude/skills/opsx-codex-apply
cp -R "$SRC/openspec-adr-generator" .claude/skills/openspec-adr-generator
```

- [ ] **Step 2: 検証**

Run: `test -f .claude/skills/opsx-codex-apply/SKILL.md && test -f .claude/skills/openspec-adr-generator/SKILL.md && find .claude/skills -type f | wc -l`
Expected: 2スキルの SKILL.md が存在し、ファイル数が元と一致（`find "$SRC" -type f | wc -l` と同数）

- [ ] **Step 3: コミット**

```bash
git add .claude/skills/
git commit -m "chore: 言語非依存のワークフロースキルを複製"
```

---

### Task 8: 最終検証

- [ ] **Step 1: ビルド・静的解析・テスト・フォーマット**

```bash
go vet ./...
go test ./...
gofmt -l backend
```
Expected: `go vet` exit 0 / `go test` 全 ok（archはskip）/ `gofmt -l` 無出力

- [ ] **Step 2: ファイル構成の目視確認**

Run: `git ls-files | sort`
Expected: `go.mod`、`AGENTS.md`、`CLAUDE.md`、`README.md`、`backend/tests/arch/*.go`、`.claude/rules/*`（7）、`.cursor/rules/*`（6）、`.codex/rules/*`（6）、`.claude/skills/**`、`docs/superpowers/**` が揃っている。

- [ ] **Step 3: 作業ツリーがクリーンであることを確認**

Run: `git status --porcelain`
Expected: 無出力（フィクスチャの消し忘れ・未追跡ファイルがない。`backend/internal` が残っていないこと）

---

## Self-Review メモ

- **Spec coverage**: spec §2成果物=Task1–7、§3層命名=全タスクで `interfaces`、§4スタック=Task3 AGENTS.md/Task5規約、§5 archテスト=Task1–2、§6規約翻訳=Task5、§7スキル=Task7、§8検証=Task8。網羅。
- **Placeholder scan**: 各ファイルは実内容を記載。`<name>`/`Xxx` はテンプレート上の意図的プレースホルダ。
- **Type consistency**: `sourceRoot`（Task1定義）をTask2が参照。`layerAfterSourceRoot` / `sourceRootName` / `forbiddenForDomain` はTask2内で完結。`expectedLayers` はTask1内で完結。
- **注意（実装者向け）**: Task5 Step1のarchitecture.mdは、説明の都合でコードフェンス内にゼロ幅文字を含めている。実ファイルでは通常の三連バッククォートで依存方向ブロックを囲み、ゼロ幅文字を含めないこと。
