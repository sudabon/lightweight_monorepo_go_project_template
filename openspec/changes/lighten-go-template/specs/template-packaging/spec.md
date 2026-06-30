## ADDED Requirements

### Requirement: 軽量版 README と AGENTS

`README.md` と `AGENTS.md` は軽量版向けに全面書き換えしなければならない（SHALL）。README には概要・ディレクトリ構成・起動方法・テスト・開発ルール・環境変数を記載する。AGENTS には人間とAIエージェント双方が読む前提で、backend の3責務ルール（handler/service/repository、handlerからDB直操作禁止）と frontend ルール、主要コマンドを簡潔に記載する。

#### Scenario: README に必須セクションが揃う

- **WHEN** `README.md` を確認する
- **THEN** 概要・ディレクトリ構成・起動方法・テスト・開発ルール・環境変数の各セクションが存在し、`go run ./cmd/server` を含む起動手順が記載されている

#### Scenario: AGENTS に責務ルールが記載される

- **WHEN** `AGENTS.md` を確認する
- **THEN** 「handler は service を呼ぶ / service は repository を呼ぶ / repository は DB を扱う / handler から DB を直接触らない」旨が記載されている

### Requirement: エージェント向けルールの簡素化

`.claude` / `.codex` / `.cursor` の `rules/` から厳密4層強制の記述を除去し、3責務ベースの軽量ルールへ簡素化しなければならない（SHALL）。重複・矛盾するアーキテクチャ規約を残してはならない（MUST NOT）。

#### Scenario: 4層強制ルールが除去される

- **WHEN** `.claude/rules/`、`.codex/rules/`、`.cursor/rules/` を確認する
- **THEN** `domain → app → interfaces → infra` の4層を強制する記述が残っておらず、handler/service/repository の3責務ルールに統一されている

### Requirement: 重厚版資産の除去

軽量化に伴い不要となる重厚版資産を除去しなければならない（SHALL）。`backend/tests/arch/dependency_rule_test.go` と `backend/tests/arch/structure_test.go` を削除し、ディレクトリ構造・依存方向をテストで強制しない。

#### Scenario: アーキテクチャテストが削除される

- **WHEN** `backend/tests/arch/` を確認する
- **THEN** `dependency_rule_test.go` と `structure_test.go` が存在しない

### Requirement: 配布用設定ファイル

PostgreSQL接続を想定した配布用設定を提供しなければならない（SHALL）。`.env.example`（`APP_ENV` / `APP_PORT` / `DATABASE_URL`）、PostgreSQLのみを含む `docker-compose.yml`、適切な `.gitignore`（`.env` を除外）を用意する。

#### Scenario: 配布用ファイルが存在する

- **WHEN** リポジトリルートを確認する
- **THEN** `.env.example`、`docker-compose.yml`、`.gitignore` が存在し、`.env.example` に `APP_ENV`・`APP_PORT`・`DATABASE_URL` が含まれる

#### Scenario: docker-compose で PostgreSQL が定義される

- **WHEN** `docker-compose.yml` を確認する
- **THEN** PostgreSQL サービスが定義され、余分なサービスを含まない

### Requirement: GitHub Template repository 化

GitHub リポジトリを Template repository として設定しなければならない（SHALL）。フォーク利用者はモジュールパス・プロジェクト名のリネームのみで開発を開始できる状態にする。

#### Scenario: リポジトリが Template として設定される

- **WHEN** GitHub 上のリポジトリ設定を確認する
- **THEN** `isTemplate` が true になっている
