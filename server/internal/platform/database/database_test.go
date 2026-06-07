package database

import (
	"context"
	"testing"

	"gorm.io/gorm"
)

func TestPingNilDB(t *testing.T) {
	err := Ping(context.Background(), nil)
	if err != gorm.ErrInvalidDB {
		t.Fatalf("err = %v, want %v", err, gorm.ErrInvalidDB)
	}
}
