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

func (r *Repository) GetComputerRMS(domain string) ([]entity.ComputerRMS, error) {
	var computers []entity.ComputerRMS
	query := `
        SELECT computer, ip
        FROM [Adman].[dbo].[UserIP]
        WHERE login LIKE '%' + $1 + '%' AND rms = 1
        GROUP BY computer, ip
    `
	_, err := r.db.Select(&computers, query, domain)
	return computers, err
}
