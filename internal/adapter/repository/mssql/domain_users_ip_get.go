package mssql

import (
	"github.com/saygik/go-userinfo/internal/entity"
)

func (r *Repository) GetDomainUsersIP(domain string) (users []entity.UserIPComputer, err error) {
	_, err = r.db.Select(&users, "GetAllUserIPByDomain $1", domain)
	return users, err
}

func (r *Repository) GetDomainUsersAvatars(domain string) (users []entity.IdNameAvatar, err error) {
	_, err = r.db.Select(&users, "GetAllUsersAvatars $1", domain)
	return users, err
}
