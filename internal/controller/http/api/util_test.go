package api

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestIsApiRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Пустая строка в конфиге не должна пропускать запросы без ключа.
	h := &Handler{allowedApi: []string{"secret-key", ""}}

	cases := []struct {
		name string
		key  string
		want bool
	}{
		{"valid key", "secret-key", true},
		{"no key", "", false},
		{"wrong key", "secret-ke", false},
		{"key from another app", "other-key", false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rec)
			c.Request = httptest.NewRequest("GET", "/", nil)
			if tc.key != "" {
				c.Request.Header.Set("X-Api-Key", tc.key)
			}

			if got := h.IsApiRequest(c); got != tc.want {
				t.Fatalf("IsApiRequest(key=%q) = %v, want %v", tc.key, got, tc.want)
			}
		})
	}
}
