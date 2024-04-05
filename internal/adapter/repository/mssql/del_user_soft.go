package mssql

import (
	"errors"

	"github.com/saygik/go-userinfo/internal/entity"
)

func (r *Repository) DelOneUserSoftware(form entity.SoftwareForm) (err error) {
	res, err := r.db.Exec("DelOneUserSoftware $1,$2", form.User, form.Id)
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
