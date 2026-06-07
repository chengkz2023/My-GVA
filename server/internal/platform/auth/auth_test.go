package auth

import (
	"context"
	"testing"
)

func TestContextWithActor(t *testing.T) {
	want := Actor{
		UserID:      1,
		AuthorityID: 888,
		Username:    "admin",
		NickName:    "Admin",
	}

	got, ok := ActorFromContext(ContextWithActor(context.Background(), want))
	if !ok {
		t.Fatal("actor not found")
	}
	if got != want {
		t.Fatalf("actor = %+v, want %+v", got, want)
	}
}

func TestActorFromContextMissing(t *testing.T) {
	_, ok := ActorFromContext(context.Background())
	if ok {
		t.Fatal("actor found, want missing")
	}
}
