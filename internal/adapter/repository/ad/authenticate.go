package ad

import "github.com/saygik/go-userinfo/internal/entity"

func (r *Repository) Authenticate(domain string, form entity.LoginForm) (bool, map[string]string, error) {
	return r.ads[domain].Authenticate(form.Email, form.Password)

}
