package usecase

import (
	"strings"
	"testing"

	"github.com/saygik/go-userinfo/internal/entity"
)

func TestBuildAdminUserPermissionDomainUserIsStandard(t *testing.T) {
	agg := &adminUserAggregate{
		upn: "user@example.local",
		roles: map[int]entity.IdNameDescription{
			1: {Id: 1, Name: "DomainUser"},
		},
		access: map[string]entity.DomainAccess{
			"example.local|user": {Domain: "example.local", AccessLevel: "user"},
			"other.local|user":   {Domain: "other.local", AccessLevel: "user"},
		},
		sections: map[string]entity.IdNameDescription{},
	}

	user, nonStandard := buildAdminUserPermission(agg, map[string]adminADUser{})

	if nonStandard {
		t.Fatalf("DomainUser with all-domain user access must be standard, reasons: %v", user.Reasons)
	}
	if user.IsSysAdmin || user.IsAdmin || user.IsTech {
		t.Fatal("DomainUser must not receive admin or tech flags")
	}
}

func TestBuildAdminUserPermissionAdminRoleIsNonStandard(t *testing.T) {
	agg := &adminUserAggregate{
		upn: "admin@example.local",
		roles: map[int]entity.IdNameDescription{
			2: {Id: 2, Name: "Admin"},
		},
		access:   map[string]entity.DomainAccess{},
		sections: map[string]entity.IdNameDescription{},
	}

	user, nonStandard := buildAdminUserPermission(agg, map[string]adminADUser{})

	if !nonStandard {
		t.Fatal("Admin role must be treated as non-standard")
	}
	if !user.IsAdmin || !user.IsTech {
		t.Fatal("Admin role must set IsAdmin and IsTech")
	}
	if !strings.Contains(strings.Join(user.Reasons, " "), "роли: Admin") {
		t.Fatalf("expected role reason, got %v", user.Reasons)
	}
}

func TestBuildAdminUserPermissionExplicitSection(t *testing.T) {
	agg := &adminUserAggregate{
		upn: "user@example.local",
		roles: map[int]entity.IdNameDescription{
			1: {Id: 1, Name: "DomainUser"},
		},
		access: map[string]entity.DomainAccess{},
		sections: map[string]entity.IdNameDescription{
			"1|/glpi/users": {Id: 1, Name: "/glpi/users"},
		},
	}

	user, nonStandard := buildAdminUserPermission(agg, map[string]adminADUser{})

	if !nonStandard {
		t.Fatal("explicit section must be treated as non-standard")
	}
	if !strings.Contains(strings.Join(user.Reasons, " "), "дополнительные разделы") {
		t.Fatalf("expected section reason, got %v", user.Reasons)
	}
}

func TestDeleteAllUserRightsGuards(t *testing.T) {
	uc := &UseCase{}

	if err := uc.DeleteAllUserRights(entity.Permissions{}, "user@example.local"); err == nil {
		t.Fatal("non-sysadmin must not delete rights")
	}

	self := entity.Permissions{IsSysAdmin: true, User: "me@example.local"}
	if err := uc.DeleteAllUserRights(self, "ME@example.local"); err == nil {
		t.Fatal("own rights must be protected")
	} else if !strings.Contains(err.Error(), "своей") {
		t.Fatalf("unexpected self-protection error: %v", err)
	}

	if err := uc.DeleteAllUserRights(self, "not-an-email"); err == nil {
		t.Fatal("invalid username must be rejected")
	}
}
