package usecase

import (
	"testing"

	"github.com/saygik/go-userinfo/internal/entity"
)

func TestHasDomainAccess(t *testing.T) {
	normalized := map[string]bool{"example.local": true}

	t.Run("case insensitive lookup", func(t *testing.T) {
		if !hasDomainAccess(normalized, "EXAMPLE.LOCAL") {
			t.Fatal("expected access for differently cased domain")
		}
	})

	t.Run("trims spaces", func(t *testing.T) {
		if !hasDomainAccess(normalized, "  Example.Local  ") {
			t.Fatal("expected access for domain with spaces")
		}
	})

	t.Run("legacy raw key", func(t *testing.T) {
		if !hasDomainAccess(map[string]bool{"EXAMPLE.LOCAL": true}, "example.local") {
			t.Fatal("expected access for raw key")
		}
	})

	t.Run("empty map", func(t *testing.T) {
		if hasDomainAccess(nil, "example.local") {
			t.Fatal("empty map must not grant access")
		}
	})
}

func TestGetAccessLevelForDomainWithMixedCase(t *testing.T) {
	uc := &UseCase{}
	perms := entity.Permissions{
		AdminDomains: map[string]bool{"example.local": true},
		TechDomains:  map[string]bool{"tech.local": true},
		UserDomains:  map[string]bool{"user.local": true},
	}

	cases := map[string]string{
		"EXAMPLE.LOCAL": "admin",
		"Tech.Local":    "tech",
		"USER.LOCAL":    "user",
		"other.local":   "none",
	}

	for domain, want := range cases {
		if got := uc.GetAccessLevelForDomain(&perms, domain); got != want {
			t.Fatalf("GetAccessLevelForDomain(%q) = %q, want %q", domain, got, want)
		}
	}
}
