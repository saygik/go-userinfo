package mssql

import (
	"errors"
)

func (r *Repository) SetOneUserSoftwareSendedToCalendar(id int64) (err error) {
	res, err := r.db.Exec("SetOneUserSoftwareSendedToCalendar $1", id)
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
