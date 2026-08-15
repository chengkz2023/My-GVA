package memory

import (
	"context"
	"sort"
	"sync"
	"time"

	"github.com/chengkz2023/My-GVA/server/internal/modules/business/example/domain"
)

// Repository 是 domain.Repository 的内存实现，示例模块无需数据库即可运行。
// 真实业务模块请参照 server/docs/how-to-add-module.md，用 GORM 实现 infrastructure/mysql。
type Repository struct {
	mu      sync.Mutex
	nextID  uint
	entries map[uint]domain.Greeting
}

func NewRepository() *Repository {
	r := &Repository{nextID: 1, entries: make(map[uint]domain.Greeting)}
	now := time.Now()
	for _, g := range []domain.Greeting{
		{Message: "你好，这是开发示例模块", Author: "example"},
		{Message: "复制本目录即可作为新模块模板", Author: "example"},
	} {
		r.entries[r.nextID] = domain.Greeting{ID: r.nextID, Message: g.Message, Author: g.Author, CreatedAt: now}
		r.nextID++
	}
	return r
}

func (r *Repository) FindByID(ctx context.Context, id uint) (domain.Greeting, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	g, ok := r.entries[id]
	if !ok {
		return domain.Greeting{}, domain.ErrGreetingNotFound
	}
	return g, nil
}

func (r *Repository) List(ctx context.Context) ([]domain.Greeting, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	items := make([]domain.Greeting, 0, len(r.entries))
	for _, g := range r.entries {
		items = append(items, g)
	}
	sort.Slice(items, func(i, j int) bool { return items[i].ID < items[j].ID })
	return items, nil
}

func (r *Repository) Create(ctx context.Context, input domain.CreateInput) (domain.Greeting, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	g := domain.Greeting{
		ID:        r.nextID,
		Message:   input.Message,
		Author:    input.Author,
		CreatedAt: time.Now(),
	}
	r.entries[r.nextID] = g
	r.nextID++
	return g, nil
}
