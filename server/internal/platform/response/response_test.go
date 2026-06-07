package response

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	apperrors "github.com/flipped-aurora/gin-vue-admin/server/internal/platform/errors"
	"github.com/gin-gonic/gin"
)

func TestErrorMapsAppError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	rec := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(rec)

	Error(ctx, apperrors.WithMessage(apperrors.NotFound, "example not found"))

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusNotFound)
	}

	var body Body
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}
	if body.Code != Failure || body.Message != "example not found" {
		t.Fatalf("body = %+v, want failure example not found", body)
	}
}

func TestErrorMapsUnknownError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	rec := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(rec)

	Error(ctx, assertErr{})

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusInternalServerError)
	}

	var body Body
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}
	if body.Code != Failure || body.Message != "internal error" {
		t.Fatalf("body = %+v, want failure internal error", body)
	}
}

type assertErr struct{}

func (assertErr) Error() string {
	return "unexpected"
}
