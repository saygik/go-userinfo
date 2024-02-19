package usecase

import (
	"github.com/saygik/go-userinfo/internal/entity"
)

func (u *UseCase) Login() error {

	return nil
}

func (u *UseCase) Authenticate(form entity.LoginForm) (bool, map[string]string, error) {
	domain := getDomainFromUserName(form.Email)
	u.ad.IsDomainExist(domain)
	if !u.ad.IsDomainExist(domain) {
		return false, nil, u.Error("нет такого домена")
	}
	return u.ad.Authenticate(domain, form)
}
