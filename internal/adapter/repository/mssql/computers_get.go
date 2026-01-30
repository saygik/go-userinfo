package mssql

import "github.com/saygik/go-userinfo/internal/entity"

// GetComputerByDomain возвращает список компьютеров домена из MSSQL
// Процедура в БД: GetComputerByDomain @domain
func (r *Repository) GetComputerByDomain(domain string) (comps []entity.DomainComputer, err error) {
	_, err = r.db.Select(&comps, "GetComputerByDomain $1", domain)
	return comps, err
}




