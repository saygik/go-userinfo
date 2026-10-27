package usecase

import (
	"reflect"
	"testing"

	"github.com/saygik/go-userinfo/internal/entity"
)

func TestNormalizeSoftUserPermissions(t *testing.T) {
	t.Run("creates main permission", func(t *testing.T) {
		got := normalizeSoftUserPermissions(nil, true)
		want := []entity.SoftUserPermission{{Name: entity.SoftUserMainPermissionName}}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("normalizeSoftUserPermissions() = %#v, want %#v", got, want)
		}
	})

	t.Run("forces main permission name", func(t *testing.T) {
		got := normalizeSoftUserPermissions([]entity.SoftUserPermission{
			{Id: 10, IdUserSoft: 20, Name: "Старое название", EndDate: "1900-01-01T00:00:00+00:00"},
			{Name: "  Доступ к отчётам  ", EndDate: "2027-01-01T00:00:00+00:00"},
		}, false)

		if got[0].Name != entity.SoftUserMainPermissionName {
			t.Fatalf("first permission name = %q, want %q", got[0].Name, entity.SoftUserMainPermissionName)
		}
		if got[1].Name != "Доступ к отчётам" {
			t.Fatalf("second permission name = %q, want %q", got[1].Name, "Доступ к отчётам")
		}
		if got[0].EndDate != "" {
			t.Fatalf("1900 date should be normalized to empty string, got %q", got[0].EndDate)
		}
		if got[0].Id != 10 || got[0].IdUserSoft != 20 {
			t.Fatalf("existing permission ids should be preserved: %#v", got[0])
		}
	})

	t.Run("clears ids when adding user", func(t *testing.T) {
		got := normalizeSoftUserPermissions([]entity.SoftUserPermission{
			{Id: 10, IdUserSoft: 20, Name: "Основные полномочия"},
		}, true)

		if got[0].Id != 0 || got[0].IdUserSoft != 0 {
			t.Fatalf("new permission ids should be cleared: %#v", got[0])
		}
	})

	t.Run("fills empty additional name", func(t *testing.T) {
		got := normalizeSoftUserPermissions([]entity.SoftUserPermission{
			{Name: entity.SoftUserMainPermissionName},
			{Name: "   "},
		}, false)

		if got[1].Name != "Дополнительные полномочия 1" {
			t.Fatalf("additional permission name = %q", got[1].Name)
		}
	})
}
