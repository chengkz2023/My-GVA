package domain

import (
	"context"
	"errors"
	"time"
)

var (
	ErrGreetingNotFound      = errors.New("greeting not found")
	ErrRepositoryUnavailable = errors.New("greeting repository unavailable")
)

// Greeting 示例领域实体：只表达业务概念，不带任何 ORM / JSON 标签。
type Greeting struct {
	ID        uint
	Message   string
	Author    string
	CreatedAt time.Time
}

// Repository 定义持久化能力，具体实现在 infrastructure 层。
// 本示例用内存实现（infrastructure/memory）；真实业务请实现 infrastructure/mysql。
type Repository interface {
	FindByID(ctx context.Context, id uint) (Greeting, error)
	List(ctx context.Context) ([]Greeting, error)
	Create(ctx context.Context, input CreateInput) (Greeting, error)
}

type CreateInput struct {
	Message string
	Author  string
}
