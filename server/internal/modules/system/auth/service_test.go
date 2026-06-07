package auth

import (
	"context"
	"testing"

	platformauth "github.com/flipped-aurora/gin-vue-admin/server/internal/platform/auth"
)

func TestServiceMe(t *testing.T) {
	service := NewService(nil)
	got, err := service.Me(platformauth.ContextWithActor(context.Background(), platformauth.Actor{
		UserID:      1,
		AuthorityID: 888,
		Username:    "admin",
		NickName:    "Admin",
	}))
	if err != nil {
		t.Fatalf("Me() error = %v", err)
	}
	if got.Actor.Username != "admin" || got.Actor.AuthorityID != 888 {
		t.Fatalf("actor = %+v, want admin authority 888", got.Actor)
	}
}

func TestServiceMeMissingActor(t *testing.T) {
	_, err := NewService(nil).Me(context.Background())
	if err == nil {
		t.Fatal("Me() error is nil, want missing actor")
	}
}
