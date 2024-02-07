package usecase

import (
	"errors"
	"strings"

	"github.com/saygik/go-userinfo/internal/entity"
)

func (u *UseCase) Login() error {

	return nil
}

func (u *UseCase) Authenticate(form entity.LoginForm) (bool, map[string]string, error) {
	domain := strings.Split(form.Email, "@")[1]
	u.ad.IsDomainExist(domain)
	if !u.ad.IsDomainExist(domain) {
		return false, nil, errors.New("нет такого домена")
	}
	return u.ad.Authenticate(domain, form)
}
