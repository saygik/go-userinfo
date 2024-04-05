package glpi

import (
	"fmt"

	"github.com/saygik/go-userinfo/internal/entity"
)

func (r *Repository) GetProblemGroups(id string) (work []entity.IdName, err error) {
	var proc = fmt.Sprintf(`
			SELECT glpi_groups.id,glpi_groups.name FROM glpi_groups_problems
			Left JOIN glpi_groups on glpi_groups_problems.groups_id=glpi_groups.id
			WHERE problems_id=%s`, id)
	_, err = r.db.Select(&work, proc)
	if err != nil {
		return nil, err
	}
	return work, nil
}
