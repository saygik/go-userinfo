package usecase

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/saygik/go-userinfo/internal/entity"
)

// GetAdminSchedules возвращает все календари вместе с выданными правами.
func (u *UseCase) GetAdminSchedules() ([]entity.Schedule, error) {
	schedules, err := u.repo.GetSchedules()
	if err != nil {
		return nil, u.Error("невозможно получить календари из базы данных")
	}

	groupNames, err := u.scheduleGroupNames()
	if err != nil {
		// Список календарей должен открываться и при недоступном GLPI:
		// тогда группы показываются по идентификатору.
		groupNames = map[int]string{}
	}

	for i := range schedules {
		users, groups, err := u.repo.GetScheduleAccessRows(schedules[i].Id)
		if err != nil {
			return nil, u.Error("невозможно получить права календаря из базы данных")
		}
		schedules[i].Managers, schedules[i].Viewers = buildScheduleAccess(users, groups, groupNames)
	}
	return schedules, nil
}

// GetScheduleGroups возвращает группы GLPI для выбора в админке.
func (u *UseCase) GetScheduleGroups() ([]entity.IdName, error) {
	groups, err := u.glpi.GetGroups()
	if err != nil {
		return nil, u.Error("невозможно получить список групп GLPI")
	}
	return groups, nil
}

// CreateSchedule создаёт календарь. Доступно только SysAdmin на уровне HTTP.
func (u *UseCase) CreateSchedule(form entity.ScheduleForm) (entity.Schedule, error) {
	if err := u.prepareScheduleForm(&form); err != nil {
		return entity.Schedule{}, err
	}
	id, err := u.repo.AddSchedule(form)
	if err != nil {
		return entity.Schedule{}, u.Error("невозможно создать календарь: " + err.Error())
	}
	return scheduleFromForm(id, form), nil
}

// UpdateSchedule изменяет календарь и полностью заменяет его права.
func (u *UseCase) UpdateSchedule(id int, form entity.ScheduleForm) error {
	if id <= 0 {
		return u.Error("не указан идентификатор календаря")
	}
	if err := u.prepareScheduleForm(&form); err != nil {
		return err
	}
	if err := u.repo.UpdateSchedule(id, form); err != nil {
		return u.Error("невозможно изменить календарь: " + err.Error())
	}
	return nil
}

// DeleteSchedule удаляет календарь, его задачи и права.
func (u *UseCase) DeleteSchedule(id int) error {
	if id <= 0 {
		return u.Error("не указан идентификатор календаря")
	}
	if err := u.repo.DeleteSchedule(id); err != nil {
		return u.Error("невозможно удалить календарь: " + err.Error())
	}
	return nil
}

func (u *UseCase) prepareScheduleForm(form *entity.ScheduleForm) error {
	form.Name = strings.TrimSpace(form.Name)
	form.Mattermost = strings.TrimSpace(form.Mattermost)
	form.Caption = strings.TrimSpace(form.Caption)

	if form.Name == "" {
		return u.Error("название календаря не заполнено")
	}
	if form.Tip != 1 && form.Tip != 2 {
		form.Tip = 1
	}

	form.Managers.Users = normalizeScheduleUsers(form.Managers.Users)
	form.Viewers.Users = normalizeScheduleUsers(form.Viewers.Users)

	groupNames := map[int]string{}
	if len(form.Managers.Groups) > 0 || len(form.Viewers.Groups) > 0 {
		var err error
		groupNames, err = u.scheduleGroupNames()
		if err != nil {
			return err
		}
	}
	var err error
	form.Managers.Groups, err = resolveScheduleGroups(form.Managers.Groups, groupNames)
	if err != nil {
		return err
	}
	form.Viewers.Groups, err = resolveScheduleGroups(form.Viewers.Groups, groupNames)
	if err != nil {
		return err
	}

	// Полный доступ важнее просмотра: не дублируем одного и того же
	// пользователя или группу в двух списках.
	form.Viewers.Users = removeUsersFromViewers(form.Viewers.Users, form.Managers.Users)
	form.Viewers.Groups = removeGroupsFromViewers(form.Viewers.Groups, form.Managers.Groups)
	return nil
}

func (u *UseCase) scheduleGroupNames() (map[int]string, error) {
	groups, err := u.glpi.GetGroups()
	if err != nil {
		return nil, u.Error("невозможно получить список групп GLPI")
	}
	names := make(map[int]string, len(groups))
	for _, group := range groups {
		names[group.Id] = group.Name
	}
	return names, nil
}

func buildScheduleAccess(users []entity.ScheduleAccessUserRow, groups []entity.ScheduleAccessGroupRow, groupNames map[int]string) (entity.ScheduleAccess, entity.ScheduleAccess) {
	managers := entity.ScheduleAccess{}
	viewers := entity.ScheduleAccess{}

	for _, row := range users {
		name := strings.TrimSpace(row.Name)
		if name == "" {
			continue
		}
		switch row.Role {
		case 1:
			managers.Users = append(managers.Users, name)
		default:
			viewers.Users = append(viewers.Users, name)
		}
	}

	for _, row := range groups {
		if row.GroupId <= 0 {
			continue
		}
		name := groupNames[row.GroupId]
		if name == "" {
			name = strconv.Itoa(row.GroupId)
		}
		group := entity.ScheduleAccessGroup{Id: row.GroupId, Name: name}
		switch row.Role {
		case 1:
			managers.Groups = append(managers.Groups, group)
		default:
			viewers.Groups = append(viewers.Groups, group)
		}
	}

	managers.Users = normalizeScheduleUsers(managers.Users)
	viewers.Users = normalizeScheduleUsers(viewers.Users)
	sortScheduleGroups(managers.Groups)
	sortScheduleGroups(viewers.Groups)
	return managers, viewers
}

func scheduleFromForm(id int, form entity.ScheduleForm) entity.Schedule {
	return entity.Schedule{
		Id:             id,
		Name:           form.Name,
		Tip:            form.Tip,
		Private:        form.Private,
		Caption:        form.Caption,
		Mattermost:     form.Mattermost,
		Managers:       form.Managers,
		Viewers:        form.Viewers,
		Available:      true,
		Edit:           true,
		ScheduleUsers:  []entity.ScheduleUser{},
		ScheduleAdmins: []entity.IdName{},
		ScheduleTasks:  []entity.ScheduleTaskCalendar{},
	}
}

func normalizeScheduleUsers(users []string) []string {
	seen := map[string]bool{}
	result := make([]string, 0, len(users))
	for _, raw := range users {
		name := strings.TrimSpace(raw)
		if name == "" {
			continue
		}
		key := strings.ToLower(name)
		if seen[key] {
			continue
		}
		seen[key] = true
		result = append(result, name)
	}
	sort.Slice(result, func(i, j int) bool {
		return strings.ToLower(result[i]) < strings.ToLower(result[j])
	})
	return result
}

func resolveScheduleGroups(groups []entity.ScheduleAccessGroup, names map[int]string) ([]entity.ScheduleAccessGroup, error) {
	seen := map[int]bool{}
	result := make([]entity.ScheduleAccessGroup, 0, len(groups))
	for _, group := range groups {
		if group.Id <= 0 || seen[group.Id] {
			continue
		}
		name, ok := names[group.Id]
		if !ok {
			return nil, fmt.Errorf("группа GLPI %d не найдена", group.Id)
		}
		seen[group.Id] = true
		result = append(result, entity.ScheduleAccessGroup{Id: group.Id, Name: name})
	}
	sortScheduleGroups(result)
	return result, nil
}

func sortScheduleGroups(groups []entity.ScheduleAccessGroup) {
	sort.Slice(groups, func(i, j int) bool {
		return strings.ToLower(groups[i].Name) < strings.ToLower(groups[j].Name)
	})
}

func removeUsersFromViewers(viewers, managers []string) []string {
	excluded := map[string]bool{}
	for _, name := range managers {
		excluded[strings.ToLower(name)] = true
	}
	result := make([]string, 0, len(viewers))
	for _, name := range viewers {
		if !excluded[strings.ToLower(name)] {
			result = append(result, name)
		}
	}
	return result
}

func removeGroupsFromViewers(viewers, managers []entity.ScheduleAccessGroup) []entity.ScheduleAccessGroup {
	excluded := map[int]bool{}
	for _, group := range managers {
		excluded[group.Id] = true
	}
	result := make([]entity.ScheduleAccessGroup, 0, len(viewers))
	for _, group := range viewers {
		if !excluded[group.Id] {
			result = append(result, group)
		}
	}
	return result
}
