package mssql

import "github.com/saygik/go-userinfo/internal/entity"

func (r *Repository) GetUserRole(userID string) (role []entity.IdNameDescription) {
	roles := []entity.IdNameDescription{}
	_, err := r.db.Select(&roles, `
		SELECT r.id, r.name, r.description
		FROM user_roles ur
		JOIN roles r ON ur.role_id = r.id
		WHERE ur.user_principal_name = $1
  		AND r.id = (
      	SELECT MIN(r2.id)
      	FROM user_roles ur2
      	JOIN roles r2 ON ur2.role_id = r2.id
      	WHERE ur2.user_principal_name = $1
  		)
		`, userID)
	if err != nil {
		roles = append(roles, entity.IdNameDescription{Id: 4, Name: "user", Description: "Пользователь"})
	}
	if len(roles) < 1 {
		roles = append(roles, entity.IdNameDescription{Id: 4, Name: "user", Description: "Пользователь"})
	}
	return roles
}
