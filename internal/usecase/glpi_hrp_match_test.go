package usecase

import (
	"testing"

	"github.com/saygik/go-userinfo/internal/entity"
)

func TestMatchSoftUsersByFIOExact(t *testing.T) {
	users := []entity.SoftUser{
		{Fio: "Иванов Иван Иванович", Login: "ivanov1"},
		{Fio: "Иванов Иван Петрович", Login: "ivanov2"},
	}
	index := buildSoftUsersByFIO(users)

	matches := matchSoftUsersByFIO(index, users, "Иванов Иван Иванович")
	if len(matches) != 1 {
		t.Fatalf("ожидалось одно точное совпадение, получено %d", len(matches))
	}
	if matches[0].Login != "ivanov1" {
		t.Fatalf("найден не тот пользователь: %s", matches[0].Login)
	}
}

func TestMatchSoftUsersByFIOFuzzyFallback(t *testing.T) {
	users := []entity.SoftUser{{Fio: "Иванов Иван Иванович", Login: "ivanov"}}
	index := buildSoftUsersByFIO(users)

	matches := matchSoftUsersByFIO(index, users, "Иванов Ивн Иванович")
	if len(matches) != 1 {
		t.Fatalf("ожидалось нечёткое совпадение, получено %d", len(matches))
	}
}

func TestMatchSoftUsersByFIOEmpty(t *testing.T) {
	users := []entity.SoftUser{{Fio: "Иванов Иван Иванович", Login: "ivanov"}}
	index := buildSoftUsersByFIO(users)

	if matches := matchSoftUsersByFIO(index, users, ""); len(matches) != 0 {
		t.Fatalf("для пустого ФИО совпадений быть не должно, получено %d", len(matches))
	}
}

func TestMatchSoftUsersByFIONoEmptyCandidateMatches(t *testing.T) {
	users := []entity.SoftUser{
		{Fio: "", Login: "without-fio"},
		{Fio: "Петров Петр Петрович", Login: "petrov"},
	}
	index := buildSoftUsersByFIO(users)

	if matches := matchSoftUsersByFIO(index, users, "Сидоров Сидор Сидорович"); len(matches) != 0 {
		t.Fatalf("запись без ФИО не должна совпадать, получено %d", len(matches))
	}
}

func TestSoftUserAlreadyDisabled(t *testing.T) {
	cases := []struct {
		name   string
		state  *entity.SoftUserAccountState
		expect bool
	}{
		{name: "нет состояния", state: nil, expect: false},
		{name: "active", state: &entity.SoftUserAccountState{Status: entity.SoftUserAccountStatusActive}, expect: false},
		{name: "disabled", state: &entity.SoftUserAccountState{Status: entity.SoftUserAccountStatusDisabled}, expect: true},
		{name: "blocked", state: &entity.SoftUserAccountState{Status: entity.SoftUserAccountStatusBlocked}, expect: true},
		{name: "deleted", state: &entity.SoftUserAccountState{Status: entity.SoftUserAccountStatusDeleted}, expect: true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			user := entity.SoftUser{AccountState: tc.state}
			if got := softUserAlreadyDisabled(user); got != tc.expect {
				t.Fatalf("ожидалось %v, получено %v", tc.expect, got)
			}
		})
	}
}
