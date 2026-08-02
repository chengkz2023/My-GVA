package auth

import (
	"context"
	"testing"

	platformauth "github.com/chengkz2023/My-GVA/server/internal/platform/auth"
)

func TestServiceMe(t *testing.T) {
	service := NewService(nil, nil, nil)
	got, err := service.Me(platformauth.ContextWithActor(context.Background(), platformauth.Actor{
		UserID:      1,
		AuthorityID: 888,
		Username:    "admin",
		NickName:    "Admin",
	}))
	if err != nil {
		t.Fatalf("Me() error = %v", err)
	}
	if got.UserInfo.NickName != "Admin" {
		t.Fatalf("userInfo = %+v, want nickName Admin", got.UserInfo)
	}
	if got.UserInfo.Authority.DefaultRouter != "authority" {
		t.Fatalf("defaultRouter = %s, want authority", got.UserInfo.Authority.DefaultRouter)
	}
}

func TestServiceMeMissingActor(t *testing.T) {
	_, err := NewService(nil, nil, nil).Me(context.Background())
	if err == nil {
		t.Fatal("Me() error is nil, want missing actor")
	}
}
