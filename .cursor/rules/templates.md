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
