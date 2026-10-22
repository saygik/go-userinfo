package oauth2authentik

import (
	"errors"
	"net/http"
	"testing"

	"github.com/saygik/go-userinfo/internal/entity"
	"golang.org/x/oauth2"
)

func TestClassifyRefreshError(t *testing.T) {
	cases := []struct {
		name     string
		err      error
		rejected bool
	}{
		{
			name:     "invalid_grant",
			err:      &oauth2.RetrieveError{ErrorCode: "invalid_grant", Response: &http.Response{StatusCode: http.StatusBadRequest, Status: "400 Bad Request"}},
			rejected: true,
		},
		{
			name:     "invalid_token",
			err:      &oauth2.RetrieveError{ErrorCode: "invalid_token", Response: &http.Response{StatusCode: http.StatusBadRequest, Status: "400 Bad Request"}},
			rejected: true,
		},
		{
			name:     "400 without error code",
			err:      &oauth2.RetrieveError{Response: &http.Response{StatusCode: http.StatusBadRequest, Status: "400 Bad Request"}},
			rejected: true,
		},
		{
			name:     "provider unavailable",
			err:      &oauth2.RetrieveError{Response: &http.Response{StatusCode: http.StatusServiceUnavailable, Status: "503 Service Unavailable"}},
			rejected: false,
		},
		{
			name:     "network error",
			err:      errors.New("dial tcp: i/o timeout"),
			rejected: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := errors.Is(classifyRefreshError(tc.err), entity.ErrRefreshTokenRejected)
			if got != tc.rejected {
				t.Fatalf("rejected = %v, want %v", got, tc.rejected)
			}
		})
	}
}
