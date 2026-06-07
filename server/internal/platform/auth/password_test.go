package auth

import "testing"

func TestBcryptPasswordHasher(t *testing.T) {
	hasher := NewBcryptPasswordHasher()

	hash, err := hasher.Hash("old-password")
	if err != nil {
		t.Fatalf("Hash() error = %v", err)
	}
	if hash == "" || hash == "old-password" {
		t.Fatalf("hash = %q, want non-empty hashed value", hash)
	}
	if !hasher.Check("old-password", hash) {
		t.Fatal("Check() = false, want true")
	}
	if hasher.Check("wrong-password", hash) {
		t.Fatal("Check() = true, want false")
	}
}
