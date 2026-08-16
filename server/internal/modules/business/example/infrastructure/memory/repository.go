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
		{Message: "你好，这是开发示例模块", Author: "example", DeptID: 1},
		{Message: "复制本目录即可作为新模块模板", Author: "example", DeptID: 2},
	} {
		r.entries[r.nextID] = domain.Greeting{ID: r.nextID, Message: g.Message, Author: g.Author, DeptID: g.DeptID, CreatedAt: now}
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

// ScopedList 行级数据权限示范（内存版）：按归属部门过滤。
// 真实业务的 mysql 实现应使用 platform/dataauth：
//
//	query := dataauth.Scope(r.db.WithContext(ctx).Model(&GreetingModel{}), "dept_id", deptIDs)
//	return query.Find(&rows).Error
//
// 其中 deptIDs 由调用方从「当前角色的数据权限映射」解析而来（见角色模块 SetDataAuthority）。
func (r *Repository) ScopedList(ctx context.Context, deptIDs []uint) ([]domain.Greeting, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	allowed := make(map[uint]bool, len(deptIDs))
	for _, id := range deptIDs {
		allowed[id] = true
	}
	items := make([]domain.Greeting, 0, len(r.entries))
	for _, g := range r.entries {
		if allowed[g.DeptID] {
			items = append(items, g)
		}
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
		DeptID:    input.DeptID,
		CreatedAt: time.Now(),
	}
	r.entries[r.nextID] = g
	r.nextID++
	return g, nil
}
