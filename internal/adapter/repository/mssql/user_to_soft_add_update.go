package mssql

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/go-gorp/gorp"
	"github.com/saygik/go-userinfo/internal/entity"
)

func nullableSoftDate(value string) interface{} {
	value = strings.TrimSpace(value)
	if value == "" || (len(value) >= 4 && (value[0:4] == "1900" || value[0:4] == "0001")) {
		return nil
	}
	return value
}

func softPermissionSended(permission entity.SoftUserPermission) bool {
	return permission.Sended && nullableSoftDate(permission.EndDate) != nil
}

func getSoftAccountStateByUserID(tx *gorp.Transaction, userID int64) (*entity.SoftUserAccountState, error) {
	var state entity.SoftUserAccountState
	err := tx.SelectOne(&state, softAccountStateSelectColumns+`
		WHERE a.id_user_software = $1`, userID)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &state, nil
}

func upsertSoftAccountState(tx *gorp.Transaction, userID int64, requested *entity.SoftUserAccountState, editor string, currentTime string) (*entity.SoftUserAccountState, error) {
	existing, err := getSoftAccountStateByUserID(tx, userID)
	if err != nil {
		return nil, err
	}

	if existing == nil {
		status := entity.SoftUserAccountStatusActive
		reason := ""
		if requested != nil {
			status = entity.NormalizeSoftUserAccountStatus(requested.Status)
			reason = strings.TrimSpace(requested.StatusReason)
		}
		if status == entity.SoftUserAccountStatusActive {
			reason = ""
		}

		_, err = tx.Exec(`
			INSERT INTO dbo.UserSoftwaresAccountStates
				(id_user_software, status, status_changed_at, status_changed_by, status_reason, source)
			VALUES ($1, $2, $3, $4, $5, 'manual')`,
			userID, status, currentTime, editor, reason)
		if err != nil {
			return nil, err
		}
		return getSoftAccountStateByUserID(tx, userID)
	}

	if requested == nil {
		return existing, nil
	}

	status := entity.NormalizeSoftUserAccountStatus(requested.Status)
	reason := strings.TrimSpace(requested.StatusReason)
	if status == entity.SoftUserAccountStatusActive {
		reason = ""
	}

	if existing.Status != status {
		_, err = tx.Exec(`
			UPDATE dbo.UserSoftwaresAccountStates
			SET status = $1,
			    status_changed_at = $2,
			    status_changed_by = $3,
			    status_reason = $4,
			    source = 'manual'
			WHERE id_user_software = $5`,
			status, currentTime, editor, reason, userID)
	} else if existing.StatusReason != reason {
		_, err = tx.Exec(`
			UPDATE dbo.UserSoftwaresAccountStates
			SET status_reason = $1,
			    source = 'manual'
			WHERE id_user_software = $2`,
			reason, userID)
	}
	if err != nil {
		return nil, err
	}

	return getSoftAccountStateByUserID(tx, userID)
}

func insertSoftPermission(tx *gorp.Transaction, userID int64, permission entity.SoftUserPermission) (entity.SoftUserPermission, error) {
	var id int64
	err := tx.QueryRow(`
		INSERT INTO dbo.UserSoftwaresPermissions (id_user_software, name, enddate, sended)
		OUTPUT INSERTED.id
		VALUES ($1, $2, $3, $4)`,
		userID, permission.Name, nullableSoftDate(permission.EndDate), softPermissionSended(permission)).Scan(&id)
	if err != nil {
		return permission, err
	}

	var result entity.SoftUserPermission
	err = tx.SelectOne(&result, softPermissionSelectColumns+`
		WHERE p.id = $1`, id)
	return result, err
}

func insertSoftPermissions(tx *gorp.Transaction, userID int64, permissions []entity.SoftUserPermission) ([]entity.SoftUserPermission, error) {
	if len(permissions) == 0 {
		permissions = []entity.SoftUserPermission{{Name: entity.SoftUserMainPermissionName}}
	}

	result := make([]entity.SoftUserPermission, 0, len(permissions))
	for _, permission := range permissions {
		inserted, err := insertSoftPermission(tx, userID, permission)
		if err != nil {
			return nil, err
		}
		result = append(result, inserted)
	}
	return result, nil
}

func getSoftPermissionsByUserID(tx *gorp.Transaction, userID int64) ([]entity.SoftUserPermission, error) {
	var permissions []entity.SoftUserPermission
	_, err := tx.Select(&permissions, softPermissionSelectColumns+`
		WHERE p.id_user_software = $1
		ORDER BY p.id`, userID)
	return permissions, err
}

func syncSoftPermissions(tx *gorp.Transaction, userID int64, permissions []entity.SoftUserPermission) ([]entity.SoftUserPermission, error) {
	if len(permissions) == 0 {
		permissions = []entity.SoftUserPermission{{Name: entity.SoftUserMainPermissionName}}
	}

	existing, err := getSoftPermissionsByUserID(tx, userID)
	if err != nil {
		return nil, err
	}

	existingByID := make(map[int64]bool, len(existing))
	for _, permission := range existing {
		existingByID[permission.Id] = true
	}

	keep := make(map[int64]bool, len(permissions))
	result := make([]entity.SoftUserPermission, 0, len(permissions))

	for _, permission := range permissions {
		if permission.Id > 0 && existingByID[permission.Id] {
			_, err = tx.Exec(`
				UPDATE dbo.UserSoftwaresPermissions
				SET name = $1, enddate = $2, sended = $3
				WHERE id = $4 AND id_user_software = $5`,
				permission.Name, nullableSoftDate(permission.EndDate), softPermissionSended(permission), permission.Id, userID)
			if err != nil {
				return nil, err
			}

			var updated entity.SoftUserPermission
			if err = tx.SelectOne(&updated, softPermissionSelectColumns+`
				WHERE p.id = $1`, permission.Id); err != nil {
				return nil, err
			}
			keep[permission.Id] = true
			result = append(result, updated)
			continue
		}

		inserted, err := insertSoftPermission(tx, userID, permission)
		if err != nil {
			return nil, err
		}
		result = append(result, inserted)
	}

	for _, permission := range existing {
		if keep[permission.Id] {
			continue
		}
		if _, err = tx.Exec(`
			DELETE FROM dbo.UserSoftwaresPermissions
			WHERE id = $1 AND id_user_software = $2`, permission.Id, userID); err != nil {
			return nil, err
		}
	}

	return result, nil
}

func (r *Repository) AddOneSoftwareUser(form entity.SoftUser, editor string, currentTime string) (formRes entity.SoftUser, err error) {
	tx, err := r.db.Begin()
	if err != nil {
		return formRes, err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	var userID int64
	err = tx.QueryRow(`
		INSERT INTO dbo.UserSoftwares
			(id_soft, user_name, user_login, user_comment, fio, [external], mail, editor, last_changes)
		OUTPUT INSERTED.id
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		form.IdSoft, form.Name, form.Login, form.Comment, form.Fio, form.External, form.Mail, editor, currentTime).Scan(&userID)
	if err != nil {
		return formRes, err
	}

	permissions, err := insertSoftPermissions(tx, userID, form.Permissions)
	if err != nil {
		return formRes, err
	}

	accountState, err := upsertSoftAccountState(tx, userID, form.AccountState, editor, currentTime)
	if err != nil {
		return formRes, err
	}

	err = tx.SelectOne(&formRes, softUserSelectColumns+`
		WHERE u.id = $1`, userID)
	if err != nil {
		return formRes, err
	}
	formRes.Permissions = permissions
	formRes.AccountState = accountState

	if err = tx.Commit(); err != nil {
		return formRes, err
	}
	return formRes, nil
}

func (r *Repository) UpdateOneSoftwareUser(form entity.SoftUser, editor string, currentTime string) (formRes entity.SoftUser, err error) {
	tx, err := r.db.Begin()
	if err != nil {
		return formRes, err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	result, err := tx.Exec(`
		UPDATE dbo.UserSoftwares
		SET user_login = $1,
		    fio = $2,
		    user_comment = $3,
		    mail = $4,
		    editor = $5,
		    last_changes = $6
		WHERE id = $7`,
		form.Login, form.Fio, form.Comment, form.Mail, editor, currentTime, form.Id)
	if err != nil {
		return formRes, err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return formRes, err
	}
	if affected < 1 {
		return formRes, fmt.Errorf("пользователь системы %d не найден", form.Id)
	}

	permissions := form.Permissions
	if permissions == nil {
		permissions, err = getSoftPermissionsByUserID(tx, form.Id)
		if err != nil {
			return formRes, err
		}
		if len(permissions) == 0 {
			permissions, err = insertSoftPermissions(tx, form.Id, nil)
			if err != nil {
				return formRes, err
			}
		}
	} else {
		permissions, err = syncSoftPermissions(tx, form.Id, permissions)
		if err != nil {
			return formRes, err
		}
	}

	accountState, err := upsertSoftAccountState(tx, form.Id, form.AccountState, editor, currentTime)
	if err != nil {
		return formRes, err
	}

	err = tx.SelectOne(&formRes, softUserSelectColumns+`
		WHERE u.id = $1`, form.Id)
	if err != nil {
		return formRes, err
	}
	formRes.Permissions = permissions
	formRes.AccountState = accountState

	if err = tx.Commit(); err != nil {
		return formRes, err
	}
	return formRes, nil
}
