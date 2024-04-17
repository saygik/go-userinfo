package glpi

import (
	"database/sql"
	"errors"
	"fmt"
)

func (r *Repository) SetHRPTicket(id int) error {
	sql := fmt.Sprintf(`INSERT INTO  glpi_tickets_hrp (tickets_id) VALUES (%d)`, id)
	res, err := r.db.Exec(sql)
	return sqlRowAffectedErrorWrapper(res, err)
}

func sqlRowAffectedErrorWrapper(res sql.Result, err error) error {
	if err != nil {
		return err
	}
	if res != nil {
		ra, err1 := res.RowsAffected()
		if err1 != nil {
			return err1
		}
		if ra < 1 {
			return errors.New("add row to sql not affected")
		}
		return nil
	}
	return nil
}
