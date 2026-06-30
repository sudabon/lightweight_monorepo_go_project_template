# 設計: monorepo_go_project_template

`monorepo_python_project_template` の **Go言語版**を本リポジトリに作成するための設計。

- 作成日: 2026-06-25
- 元テンプレート: `/Users/suda/workspace/sudabon/monorepo_python_project_template`
- スコープ: **スケルトンのみ**（アーキテクチャテスト + エージェント規約 + ドキュメント）。アプリ本体コードは作成しない（Python版と同じ思想）。

## 1. 目的・方針

Python版テンプレートは「Clean Architectureの依存方向とソース配置をテストで固定し、エージェント向けの規約をリポジトリ内にまとめた、実装着手前のベース」である。Go版もこの思想を忠実に踏襲する。

- 実アプリ（`backend/internal/` 配下のパッケージ、フロントエンド本体）は**未配置**。
- 技術スタックは未実装でも「採用予定」として `AGENTS.md`・規約ドキュメントに明記する（Python版がFastAPI/SQLAlchemy等を未実装でも記載しているのと同様）。
- アーキテクチャテストは対象パッケージ追加後に有効化される（素の状態では `t.Skip()` で緑）。

## 2. 成果物（ファイル構成）

Python版と1:1対応するスケルトン構成。

```
├── README.md                      # Go版テンプレート説明（書き換え）
├── AGENTS.md                      # Goプロジェクト概要・スタック・コマンド（書き換え）
├── CLAUDE.md                      # @AGENTS.md（同一）
├── go.mod                         # モジュール定義（archテストのコンパイル/実行に必要な最小要件）
├── .claude/rules/                 # Go向け規約（architecture/backend-conventions/frontend-conventions/testing/git/templates + codex-guidelines）
├── .cursor/rules/                 # 同内容（codex-guidelines.md を除く6ファイル）
├── .codex/rules/                  # 同内容（codex-guidelines.md を除く6ファイル）
├── .claude/skills/                # opsx-codex-apply / openspec-adr-generator を**そのまま複製**（言語非依存のワークフロー）
└── backend/tests/arch/            # Go版アーキテクチャテスト
    ├── dependency_rule_test.go
    └── structure_test.go
```

- 既存の `.gitignore`（Go用）・`LICENSE` はそのまま使用。
- `.uv-cache/`（未追跡のゴミ）は変更しない。

## 3. 層命名の決定（Go予約語対応）

Goでは `interface` が予約語のためパッケージ名に使用できない。Python版の `domain / app / interface / infra` を次の通り対応させる。

| Python版 | Go版 | 理由 |
|----------|------|------|
| domain | `domain` | 同一 |
| app | `app` | 同一 |
| interface | **`interfaces`**（複数形） | 予約語回避・Python版との見た目の差を最小化 |
| infra | `infra` | 同一 |

依存方向（厳守）:

```
interfaces → app → domain
infra      → app → domain
```

- domain は他のどの層にも依存しない（標準ライブラリ中心）。
- app は domain のみに依存する。
- infra は app のリポジトリ/portインターフェースを実装する。
- interfaces は app のusecaseを呼び出す。

## 4. 技術スタックのマッピング

`AGENTS.md` および各規約ドキュメントに記載する内容。

| 項目 | Python版 | Go版 |
|------|----------|------|
| 言語 | Python 3.13 | Go 1.23+ |
| Webフレームワーク | FastAPI | Echo（v4） |
| パッケージ管理 | uv | go modules |
| ORM/DB | SQLAlchemy / PostgreSQL | GORM v2 / PostgreSQL |
| マイグレーション | Alembic | golang-migrate（SQLファイル） |
| テスト | pytest | `go test`（table-driven） |
| リンター | Ruff | golangci-lint + gofmt |
| 型チェック | mypy | （静的型付けのため不要、`go vet`） |
| DI | FastAPI `Depends` | コンストラクタ注入（手動配線。必要なら google/wire を言及） |
| 設定 | pydantic-settings | 環境変数由来（`infra/config`） |
| ロギング | （記載なし） | `log/slog`（`infra/logging`） |

フロントエンド規約（React / TypeScript / Vite / Shadcn/ui / TanStack Query）は**変更なしで維持**する（言語非依存のため）。

### 想定ディレクトリ構成（AGENTS.md に記載、本タスクでは未作成）

```
backend/
├── cmd/
│   └── server/
│       └── main.go            # サーバプロセス管理（エントリポイント）
├── internal/
│   └── todo_app/              # ソースルート（go.mod のモジュールパス末尾と整合。フォーク時にリネーム）
│       ├── app/               # Application層
│       │   ├── common/        # 結果/エラー型
│       │   ├── dto/           # 層間データ受け渡し
│       │   ├── repository/    # 永続化インターフェース定義
│       │   ├── port/          # 外部サービスのインターフェース定義
│       │   └── usecase/       # ビジネスロジック実行フロー
│       ├── domain/            # Domain層
│       │   ├── entity/        # エンティティ
│       │   └── service/       # ドメインサービス
│       ├── infra/             # Infrastructure層
│       │   ├── cli/           # CLIコマンド
│       │   ├── config/        # 設定（環境変数）
│       │   ├── notification/  # 通知
│       │   ├── logging/       # ロギング（slog）
│       │   ├── persistence/   # 永続化（GORMモデル・リポジトリ実装）
│       │   └── web/           # Echoルーター・ミドルウェア
│       └── interfaces/        # Interface層
│           ├── controller/    # 入力リクエスト処理
│           ├── presenter/     # 出力変換
│           └── viewmodel/     # レスポンス整形
├── migrations/                # golang-migrate のSQLファイル
└── tests/
    └── arch/                  # アーキテクチャテスト（本タスクで作成）
```

## 5. アーキテクチャテストの方針（Go実装）

Python版の `ast` 解析を、Go標準の `go/parser` + `go/ast` で再現する（標準ライブラリのみ、外部依存なし）。

### dependency_rule_test.go

- ソースルート（`backend/internal/todo_app/`）配下の `domain/` 内の全 `.go` ファイルのimportを走査する。
- import先のパスに `/app`・`/interfaces`・`/infra` 層が含まれる場合を違反として収集する。
- 違反があればテスト失敗、メッセージに違反箇所を列挙する。

### structure_test.go

- ソースルート直下のディレクトリ集合が、期待する4層 `domain / app / interfaces / infra` のみであることを検査する。
- 必須層の欠落・想定外フォルダの存在を失敗とする。

### スケルトン挙動

- 対象ソースルートが**存在しない場合は `t.Skip()`** で緑にする。
- これにより素のテンプレートで `go test ./...` がパスし、検証可能になる。
- Python版README同様、「パッケージ追加後に有効化される」旨を README / テストコメントに明記する。

### モジュールパス・ソースルートの扱い

- `go.mod` のモジュールパスは `github.com/sudabon/todo_app`（フォーク時にリネームする旨を注記）。
- アーキテクチャテストはソースルートを**ファイルシステム上の相対パス**（テストファイルからの相対）で解決し、モジュールパスに依存しない（対象パッケージをimportしないためテスト単体でコンパイル可能）。

## 6. 各規約ドキュメントの翻訳方針

| ファイル | 方針 |
|----------|------|
| architecture.md | 依存方向を `interfaces → app → domain` / `infra → app → domain` に。禁止事項をGo文脈（domain層でGORM/Echoをimportしない、usecaseで具象を直接生成しない、controllerからentityを直接返さない 等）に翻訳。 |
| backend-conventions.md | Go規約へ翻訳: 公開要素はdoc comment、命名（公開PascalCase/非公開camelCase/パッケージ名小文字）、`any`/`interface{}` の濫用禁止、`context.Context` の伝播、エラーは値で返す。層別ルール（entityはstruct+コンストラクタでバリデーション、usecaseは単一メソッド`Execute`、repository/portはGoのinterface、GORMモデルはinfra/persistenceでdomain entityと分離し変換メソッドを用意、設定は`infra/config`、ロギングは`infra/logging`でslog、Echoハンドラは薄く保ちcontrollerへ委譲、DIはコンストラクタ注入）。 |
| frontend-conventions.md | **変更なしで維持**（言語非依存）。 |
| testing.md | Go testへ翻訳: table-driven、`go test ./...`、モックはinterface実装（手書き or gomock/mockery）、統合テストはbuild tag/testcontainers、命名 `TestXxx_条件_期待結果`。フロント部分は維持。 |
| git.md | **変更なしで維持**（日本語コミットメッセージ・プレフィックス・1コミット1関心事・.env非コミット）。 |
| templates.md | Goコードテンプレートへ翻訳: entity（struct + コンストラクタ）、usecase（struct + `Execute`）、repositoryインターフェース、Echoルーター/ハンドラ、portインターフェース。 |
| codex-guidelines.md | **変更なしで維持**（`.claude/rules/` のみ）。 |

## 7. スキルの扱い

`.claude/skills/opsx-codex-apply/` および `.claude/skills/openspec-adr-generator/` は言語非依存のワークフローのため、**ファイル内容をそのまま複製**する（`generate_adr.py` を含むスクリプトも変更しない）。

## 8. 検証方法

- `go build ./...` … 成功すること。
- `go vet ./...` … 警告なし。
- `go test ./...` … 緑（スケルトンのためarchテストはskip）。
- `gofmt -l .` … 差分なし（Goファイルが整形済み）。
- ファイル構成がPython版と1:1対応していることを目視確認。

## 9. 非対象（YAGNI）

- バックエンドのアプリ本体コード（`internal/` 配下の実パッケージ）。
- フロントエンドの実コード。
- CI設定・Dockerfile・docker-compose（Python版テンプレートにも含まれないため）。
- `golangci-lint` の設定ファイル（必要なら別タスク）。
