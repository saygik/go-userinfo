package usecase

import (
	"strconv"
	"strings"

	"github.com/saygik/go-userinfo/internal/entity"
)

const defaultSoftwareMailSubject = "Срок действия вашей учетной записи истекает"

const defaultSoftwareMailBody = `Срок действия полномочий «{{permission}}» учетной записи {{login}} пользователя {{fio}} для системы {{system}} истекает {{enddate}}.
Примите меры для перерегистрации в системе или пропустите это письмо, если данная система вас не интересует.`

// defaultSoftwareMailTemplate возвращает общий шаблон письма об истечении
// полномочий. Он используется, если для системы не задан собственный шаблон.
func (u *UseCase) defaultSoftwareMailTemplate() entity.SoftwareMailTemplate {
	return entity.SoftwareMailTemplate{
		Subject:        defaultSoftwareMailSubject,
		Body:           defaultSoftwareMailBody,
		DefaultSubject: defaultSoftwareMailSubject,
		DefaultBody:    defaultSoftwareMailBody,
	}
}

// softwareMailTemplate возвращает действующий шаблон системы: пользовательский,
// если он сохранён, иначе общий. Ошибка чтения БД не должна ломать карточку
// системы или рассылку — в этом случае остаётся общий шаблон.
func (u *UseCase) softwareMailTemplate(idSoft int64) entity.SoftwareMailTemplate {
	template := u.defaultSoftwareMailTemplate()

	custom, found, err := u.repo.GetSoftwareMailTemplate(idSoft)
	if err != nil {
		if u.log != nil {
			u.log.Errorf("не удалось получить шаблон письма системы %d: %v", idSoft, err)
		}
		return template
	}
	if !found {
		return template
	}

	template.Subject = custom.Subject
	template.Body = custom.Body
	template.IsCustom = true
	template.UpdatedAt = custom.UpdatedAt
	template.UpdatedBy = custom.UpdatedBy
	return template
}

// renderSoftwareMailTemplate подставляет значения в шаблон. Неизвестные
// placeholder-ы остаются без изменений, чтобы ошибка в тексте была заметна.
func renderSoftwareMailTemplate(text string, values map[string]string) string {
	replacements := make([]string, 0, len(values)*2)
	for key, value := range values {
		replacements = append(replacements, "{{"+key+"}}", value)
	}
	return strings.NewReplacer(replacements...).Replace(text)
}

func softwareMailTemplateValues(permission entity.SoftUserPermissionEOL, softName string) map[string]string {
	return map[string]string{
		"permission": permission.PermissionName,
		"login":      permission.Login,
		"fio":        permission.Fio,
		"system":     softName,
		"enddate":    parseDate(permission.EndDate),
	}
}

// SaveSoftwareMailTemplate сохраняет шаблон письма для системы. Изменять
// шаблон может SysAdmin или администратор этой системы.
func (u *UseCase) SaveSoftwareMailTemplate(ids string, form entity.SoftwareMailTemplate, currentUser string, perms entity.Permissions) (entity.SoftwareMailTemplate, error) {
	id, err := u.softwareMailTemplateID(ids, currentUser, perms)
	if err != nil {
		return entity.SoftwareMailTemplate{}, err
	}

	subject := strings.TrimSpace(form.Subject)
	body := strings.TrimSpace(form.Body)
	if subject == "" {
		return entity.SoftwareMailTemplate{}, u.Error("тема письма не заполнена")
	}
	if body == "" {
		return entity.SoftwareMailTemplate{}, u.Error("текст письма не заполнен")
	}
	if len([]rune(subject)) > 300 {
		return entity.SoftwareMailTemplate{}, u.Error("тема письма не должна превышать 300 символов")
	}

	if err := u.repo.SaveSoftwareMailTemplate(int64(id), subject, body, currentUser); err != nil {
		return entity.SoftwareMailTemplate{}, u.Error("невозможно сохранить шаблон письма: " + err.Error())
	}
	return u.softwareMailTemplate(int64(id)), nil
}

// DeleteSoftwareMailTemplate удаляет пользовательский шаблон системы,
// возвращая общий шаблон по умолчанию.
func (u *UseCase) DeleteSoftwareMailTemplate(ids string, currentUser string, perms entity.Permissions) (entity.SoftwareMailTemplate, error) {
	id, err := u.softwareMailTemplateID(ids, currentUser, perms)
	if err != nil {
		return entity.SoftwareMailTemplate{}, err
	}

	if err := u.repo.DeleteSoftwareMailTemplate(int64(id)); err != nil {
		return entity.SoftwareMailTemplate{}, u.Error("невозможно сбросить шаблон письма: " + err.Error())
	}
	return u.softwareMailTemplate(int64(id)), nil
}

// softwareMailTemplateID проверяет идентификатор системы и права на изменение
// её шаблона.
func (u *UseCase) softwareMailTemplateID(ids string, currentUser string, perms entity.Permissions) (int, error) {
	id, err := strconv.Atoi(ids)
	if err != nil || id <= 0 {
		return 0, u.Error("неправильный id системы")
	}
	if perms.IsSysAdmin {
		return id, nil
	}

	software, err := u.GetSoftware(ids)
	if err != nil {
		return 0, err
	}
	for _, admin := range software.Admins {
		if softwareAdminMatchesUser(admin, currentUser) {
			return id, nil
		}
	}
	return 0, u.Error("нет прав на изменение шаблона письма этой системы")
}

// softwareAdminMatchesUser сравнивает текущего пользователя с администратором
// системы по логину (name) и UPN. Регистр не важен.
func softwareAdminMatchesUser(admin map[string]interface{}, currentUser string) bool {
	currentUser = strings.TrimSpace(currentUser)
	if currentUser == "" {
		return false
	}
	for _, key := range []string{"name", "userPrincipalName", "upn"} {
		value, ok := admin[key].(string)
		if ok && strings.EqualFold(strings.TrimSpace(value), currentUser) {
			return true
		}
	}
	return false
}
