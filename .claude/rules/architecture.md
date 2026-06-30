# アーキテクチャ原則

backend は `handler / service / repository` の3責務を基本にする。

## 責務

- `handler`: Echo の HTTP 入出力、リクエストの取り出し、レスポンス整形。
- `service`: 業務判断、ユースケースの流れ、必要な repository 呼び出し。
- `repository`: DB・SQL・ドライバ固有処理。
- `config`: 環境変数読み込み。
- `db`: DB接続初期化。
- `router`: Echo ルート登録と依存の配線。

## 禁止事項

- `handler` に SQL や DB ドライバ呼び出しを書かない。
- `handler` から `repository` や `db` を直接呼ばない。
- `service` に HTTP 固有の `echo.Context` を渡さない。
- `repository` から HTTP レスポンスを組み立てない。
- `os.Getenv` は `backend/internal/config` 以外で呼ばない。
