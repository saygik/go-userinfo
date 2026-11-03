package usecase

import (
	"testing"

	"github.com/saygik/go-userinfo/internal/entity"
)

type stubPublicUsersRedis struct {
	Redis
	users map[string]string
}

func (s stubPublicUsersRedis) GetKeyFieldAll(string) (map[string]string, error) {
	return s.users, nil
}

func TestGetADUsersPublicInfoExcludesServiceAccounts(t *testing.T) {
	uc := &UseCase{
		redis: stubPublicUsersRedis{
			users: map[string]string{
				"normal":   `{"userPrincipalName":"normal@example.local","cn":"Normal User","service":false,"domain":{"name":"example.local"}}`,
				"service":  `{"userPrincipalName":"service@example.local","cn":"Service User","service":true,"domain":{"name":"example.local"}}`,
				"disabled": `{"userPrincipalName":"disabled@example.local","cn":"Disabled User","disabled":true,"domain":{"name":"example.local"}}`,
			},
		},
	}

	// Даже системный администратор не должен получать служебные учётные записи
	// через публичный список.
	users, err := uc.GetADUsersPublicInfo(entity.Permissions{IsTech: true, IsSysAdmin: true})
	if err != nil {
		t.Fatalf("GetADUsersPublicInfo() error = %v", err)
	}

	if len(users) != 1 {
		t.Fatalf("GetADUsersPublicInfo() returned %d users, want 1", len(users))
	}

	if got := users[0]["name"]; got != "normal@example.local" {
		t.Fatalf("unexpected user in public list: %v", got)
	}
}
