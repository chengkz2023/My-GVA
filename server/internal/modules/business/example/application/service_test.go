package application

import (
	"context"
	"errors"
	"testing"

	"github.com/chengkz2023/My-GVA/server/internal/modules/business/example/domain"
	apperrors "github.com/chengkz2023/My-GVA/server/internal/platform/errors"
)

func TestGetNotFoundMapsToNotFound(t *testing.T) {
	service := NewService(fakeRepo{getErr: domain.ErrGreetingNotFound})
	_, err := service.Get(context.Background(), 99)
	var appErr *apperrors.Error
	if !errors.As(err, &appErr) || appErr.Kind != apperrors.NotFound {
		t.Fatalf("err = %v, want not found app error", err)
	}
}

func TestListGreetings(t *testing.T) {
	service := NewService(fakeRepo{items: []domain.Greeting{{ID: 1, Message: "hi", Author: "a"}}})
	got, err := service.List(context.Background())
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if len(got.List) != 1 || got.List[0].Message != "hi" {
		t.Fatalf("list = %+v, want one greeting", got.List)
	}
}

func TestCreateRequiresMessage(t *testing.T) {
	service := NewService(fakeRepo{})
	_, err := service.Create(context.Background(), CreateGreetingCommand{})
	var appErr *apperrors.Error
	if !errors.As(err, &appErr) || appErr.Kind != apperrors.Validation {
		t.Fatalf("err = %v, want validation app error", err)
	}
}

func TestCreateGreeting(t *testing.T) {
	service := NewService(fakeRepo{})
	got, err := service.Create(context.Background(), CreateGreetingCommand{Message: "hello", Author: "me"})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if got.Greeting.Message != "hello" || got.Greeting.Author != "me" {
		t.Fatalf("greeting = %+v, want hello from me", got.Greeting)
	}
}

type fakeRepo struct {
	items  []domain.Greeting
	getErr error
}

func (f fakeRepo) FindByID(ctx context.Context, id uint) (domain.Greeting, error) {
	if f.getErr != nil {
		return domain.Greeting{}, f.getErr
	}
	for _, g := range f.items {
		if g.ID == id {
			return g, nil
		}
	}
	return domain.Greeting{}, domain.ErrGreetingNotFound
}

func (f fakeRepo) List(ctx context.Context) ([]domain.Greeting, error) {
	return f.items, nil
}

func (f fakeRepo) Create(ctx context.Context, input domain.CreateInput) (domain.Greeting, error) {
	return domain.Greeting{ID: 1, Message: input.Message, Author: input.Author}, nil
}
