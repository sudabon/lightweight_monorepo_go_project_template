# ファイル作成時のテンプレート

## 新規 service

```go
package service

type XxxService struct {
	repo XxxRepository
}

func NewXxxService(repo XxxRepository) *XxxService {
	return &XxxService{repo: repo}
}

func (s *XxxService) Run(input XxxInput) (XxxOutput, error) {
	return XxxOutput{}, nil
}
```

## 新規 handler

```go
package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type XxxService interface {
	Run(input XxxInput) (XxxOutput, error)
}

type XxxHandler struct {
	service XxxService
}

func NewXxxHandler(service XxxService) *XxxHandler {
	return &XxxHandler{service: service}
}

func (h *XxxHandler) Get(c echo.Context) error {
	output, err := h.service.Run(XxxInput{})
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, output)
}
```

## 新規 repository

```go
package repository

import (
	"context"
	"database/sql"
)

type XxxRepository struct {
	db *sql.DB
}

func NewXxxRepository(db *sql.DB) *XxxRepository {
	return &XxxRepository{db: db}
}

func (r *XxxRepository) Find(ctx context.Context, id string) (XxxRecord, error) {
	return XxxRecord{}, nil
}
```
