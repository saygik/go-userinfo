package mssql

import "github.com/saygik/go-userinfo/internal/entity"

func (r *Repository) GetOneDelegateMailBoxes(fio string) (boxes []entity.MailBoxDelegates, err error) {
	_, err = r.db.Select(&boxes, "GetOneDelegateMailBoxes $1", fio)
	return boxes, err
}
