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
	searchPattern := "%" + domain + "%"
	query := `
        WITH LastRecords AS (
            SELECT computer, ip, rms, date,
                ROW_NUMBER() OVER (
                    PARTITION BY computer, ip
                    ORDER BY date DESC
                ) as rn
            FROM [Adman].[dbo].[UserIP]
            WHERE login LIKE $1
        )
        SELECT computer, ip, rms
        FROM LastRecords
        WHERE rn = 1
    `
	_, err := r.db.Select(&computers, query, searchPattern)
	return computers, err
}
