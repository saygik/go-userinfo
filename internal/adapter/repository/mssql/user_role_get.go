package mssql

import "github.com/saygik/go-userinfo/internal/entity"

func (r *Repository) GetUserRole(userID string) (role entity.IdName) {
	roles := []entity.IdName{}
	_, err := r.db.Select(&roles, "GetUserRoles $1", userID)
	if err != nil {
		return entity.IdName{Id: 6, Name: "Пользователь"}
	}
	if len(roles) < 1 {
		return entity.IdName{Id: 6, Name: "Пользователь"}
	}
	return roles[0]
}
