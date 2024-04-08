package mssql

import "github.com/saygik/go-userinfo/internal/entity"

func (r *Repository) UpdateScheduleTask(form entity.ScheduleTask) error {
	res, err := r.db.Exec("UpdateScheduleTask $1,$2,$3,$4,$5,$6,$7,$8,$9", form.Id, form.Tip, form.Status, form.Title, form.Start, form.End, form.AllDay, form.SendMattermost, form.Comment)
	return sqlRowAffectedErrorWrapper(res, err)
}
