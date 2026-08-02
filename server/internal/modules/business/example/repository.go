package example

import "context"

type Repository interface {
	Info(ctx context.Context) Info
}

type MemoryRepository struct{}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{}
}

func (r *MemoryRepository) Info(ctx context.Context) Info {
	return Info{
		Name:    "business/example",
		Message: "module registered",
	}
}
