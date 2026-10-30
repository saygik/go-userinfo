package usecase

import (
	"testing"

	"github.com/saygik/go-userinfo/internal/entity"
)

func TestIsServiceAccountUser(t *testing.T) {
	groups := []string{
		"CN=Service Accounts,OU=Groups,DC=example,DC=local",
		"CN=Monitoring,OU=Groups,DC=example,DC=local",
	}

	t.Run("matches group CN from DN", func(t *testing.T) {
		if !isServiceAccountUser([]string{"Domain Users", "Service Accounts"}, groups) {
			t.Fatal("expected user to be a service account")
		}
	})

	t.Run("matching is case insensitive", func(t *testing.T) {
		if !isServiceAccountUser([]string{"service accounts"}, groups) {
			t.Fatal("expected case insensitive match")
		}
	})

	t.Run("ignores regular groups", func(t *testing.T) {
		if isServiceAccountUser([]string{"Domain Users", "VPN"}, groups) {
			t.Fatal("expected user not to be a service account")
		}
	})

	t.Run("empty configuration", func(t *testing.T) {
		if isServiceAccountUser([]string{"Service Accounts"}, nil) {
			t.Fatal("no configured groups must not produce service accounts")
		}
	})

	t.Run("plain group name in configuration", func(t *testing.T) {
		if !isServiceAccountUser([]string{"Monitoring"}, []string{" Monitoring "}) {
			t.Fatal("expected plain group name to match")
		}
	})

	t.Run("escaped comma in DN", func(t *testing.T) {
		if !isServiceAccountUser([]string{"Service, Accounts"}, []string{"CN=Service\\, Accounts,OU=Groups,DC=example,DC=local"}) {
			t.Fatal("expected escaped comma CN to match")
		}
	})
}

func TestIsServiceMailUser(t *testing.T) {
	t.Run("matches plain group name", func(t *testing.T) {
		if !isServiceMailUser([]string{"Domain Users", "Служебная почта"}) {
			t.Fatal("expected user to have service mail")
		}
	})

	t.Run("matching is case insensitive and trims spaces", func(t *testing.T) {
		if !isServiceMailUser([]string{"  служебная ПОЧТА  "}) {
			t.Fatal("expected case insensitive match")
		}
	})

	t.Run("matches CN from DN", func(t *testing.T) {
		if !isServiceMailUser([]string{"CN=Служебная почта,OU=Groups,DC=example,DC=local"}) {
			t.Fatal("expected DN match")
		}
	})

	t.Run("ignores other groups", func(t *testing.T) {
		if isServiceMailUser([]string{"Domain Users", "Service Accounts"}) {
			t.Fatal("expected user not to have service mail")
		}
	})

	t.Run("empty member list", func(t *testing.T) {
		if isServiceMailUser(nil) {
			t.Fatal("empty member list must not match")
		}
	})
}

func TestCanSeeServiceAccounts(t *testing.T) {
	sysadmin := entity.Permissions{IsSysAdmin: true}
	if !(&UseCase{}).canSeeServiceAccounts(&sysadmin, "example.local") {
		t.Fatal("sysadmin must see service accounts in any domain")
	}

	domainAdmin := entity.Permissions{AdminDomains: map[string]bool{"example.local": true}}
	if !(&UseCase{}).canSeeServiceAccounts(&domainAdmin, "example.local") {
		t.Fatal("domain admin must see service accounts in own domain")
	}
	if !(&UseCase{}).canSeeServiceAccounts(&domainAdmin, "EXAMPLE.LOCAL") {
		t.Fatal("domain admin must see service accounts in own domain regardless of case")
	}
	if (&UseCase{}).canSeeServiceAccounts(&domainAdmin, "other.local") {
		t.Fatal("domain admin must not see service accounts in other domains")
	}

	tech := entity.Permissions{TechDomains: map[string]bool{"example.local": true}}
	if (&UseCase{}).canSeeServiceAccounts(&tech, "example.local") {
		t.Fatal("tech must not see service accounts")
	}

	if (&UseCase{}).canSeeServiceAccounts(nil, "example.local") {
		t.Fatal("nil permissions must not see service accounts")
	}
}

func TestUserDomainName(t *testing.T) {
	user := map[string]interface{}{
		"domain": map[string]interface{}{"name": "example.local", "title": "Example"},
	}
	if got := userDomainName(user); got != "example.local" {
		t.Fatalf("userDomainName() = %q, want %q", got, "example.local")
	}
	if got := userDomainName(map[string]interface{}{"domain": "plain.local"}); got != "plain.local" {
		t.Fatalf("userDomainName() = %q, want %q", got, "plain.local")
	}
	if got := userDomainName(map[string]interface{}{}); got != "" {
		t.Fatalf("userDomainName() = %q, want empty string", got)
	}
}
