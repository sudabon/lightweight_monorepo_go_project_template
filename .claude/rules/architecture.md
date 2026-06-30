# アーキテクチャ原則

## 依存の方向（厳守）

```
interfaces → app → domain
infra      → app → domain
```

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
