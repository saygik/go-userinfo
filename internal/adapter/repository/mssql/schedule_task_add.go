package mssql

import "github.com/saygik/go-userinfo/internal/entity"

func (r *Repository) AddScheduleTask(form entity.ScheduleTask) (formRes entity.ScheduleTask, err error) {
	err = r.db.QueryRow("AddScheduleTask $1,$2,$3,$4,$5,$6,$7,$8,$9,$10",
		form.Idc, form.Tip, form.Status, form.Title, form.Upn, form.Start, form.End, form.AllDay, form.SendMattermost, form.Comment).Scan(&formRes.Id,
		&formRes.Idc, &formRes.Tip, &formRes.Status, &formRes.Title, &formRes.Upn, &formRes.Start, &formRes.End, &formRes.AllDay, &formRes.SendMattermost, &formRes.Comment)

	return formRes, err
}
