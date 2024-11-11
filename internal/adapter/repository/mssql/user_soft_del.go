package mssql

import (
	"errors"
)

func (r *Repository) DelOneUserSoftware(id string) (err error) {
	res, err := r.db.Exec("DelOneUserSoftware $1", id)
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
