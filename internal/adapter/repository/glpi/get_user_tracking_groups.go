package glpi

import (
	"fmt"

	"github.com/saygik/go-userinfo/internal/entity"
)

// ************************* Незакрытые заявки **********************************//

func (r *Repository) GetUserTrackingGroups(ids string) (groups []entity.IdName, err error) {
	sql := fmt.Sprintf(`SELECT id, name  from glpi_groups WHERE id IN (%s)`, ids)
	_, err = r.db.Select(&groups, sql)
	return groups, err
}
