package usecase

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/saygik/go-userinfo/internal/entity"
)

const (
	// Error messages
	errDomainNotFound   = "домен не найден"
	errUserNotFound     = "пользователь не существует"
	errUserDNNotDefined = "не орпеделен dn пользователя"
	errNoPermission     = "У Вас нет прав на эту операцию"
	errUserNotInSystem  = "нет таккого пользователя в системе"
	//	errCashAllNotFound    = "нет общего кеша для обновления"
	errInvalidGroup       = "Неправильная группа: %s. Должна быть 'whitelist', 'full', 'tech' or ''"
	errFailedToAddToGroup = "failed to add to group %s: %w"

	// Group name constants
	groupNameWhitelist = "whitelist"
	groupNameFull      = "full"
	groupNameTech      = "tech"

	// AD error indicators
	adErrorAccessDenied  = "INSUFF_ACCESS_RIGHTS"
	adErrorEntryExists   = "ENTRY_EXISTS"
	adErrorAlreadyExists = "уже существует"
	adErrorNotExists     = "00000561"
	adErrorNotExistsText = "не существует"

	// Permission group
	permissionGroupDomainAdmins = "Администраторы пользователей домена"

	// Redis keys
	redisKeyAllUsers        = "allusers"
	redisKeyTempGroupChange = "temp_group_change"

	// User fields
	fieldDistinguishedName = "distinguishedName"
	fieldcn                = "cn"
	fieldMemberOf          = "memberOf"
	fieldUserPrincipalName = "userPrincipalName"
	fieldInternetWL        = "internetwl"
	fieldInternet          = "internet"
	fieldInternetTech      = "internettech"
)

// getUserDN extracts the distinguished name from user info map
func getUserDN(userInfo map[string]any) (string, error) {
	dn, ok := userInfo[fieldDistinguishedName].(string)
	if !ok || dn == "" {
		return "", errors.New(errUserDNNotDefined)
	}
	return dn, nil
}

func getUsercn(userInfo map[string]any) (string, error) {
	dn, ok := userInfo[fieldcn].(string)
	if !ok || dn == "" {
		return "", errors.New(errUserDNNotDefined)
	}
	return dn, nil
}

// validateDomainAndGetUser validates domain existence and retrieves user information
func (u *UseCase) validateDomainAndGetUser(user string) (string, map[string]any, error) {
	domain := getDomainFromUserName(user)
	if !u.ad.IsDomainExist(domain) {
		return "", nil, u.Error(errDomainNotFound)
	}

	userInfo, err := u.ad.GetUser(domain, user)
	if err != nil {
		return "", nil, u.Error(errUserNotFound)
	}

	return domain, userInfo, nil
}

// isEntryExistsError checks if the error indicates that an entry already exists
func isEntryAccessError(err error) bool {
	if err == nil {
		return false
	}
	errMsg := err.Error()
	return strings.Contains(errMsg, adErrorAccessDenied)
}

// isEntryExistsError checks if the error indicates that an entry already exists
func isEntryExistsError(err error) bool {
	if err == nil {
		return false
	}
	errMsg := err.Error()
	return strings.Contains(errMsg, adErrorEntryExists) || strings.Contains(errMsg, adErrorAlreadyExists)
}

// isEntryNotExistsError checks if the error indicates that an entry does not exist
func isEntryNotExistsError(err error) bool {
	if err == nil {
		return false
	}
	errMsg := err.Error()
	return strings.Contains(errMsg, adErrorNotExists) || strings.Contains(errMsg, adErrorNotExistsText)
}

func (u *UseCase) ADUserAddGroup(user string, group string) error {
	domain, userInfo, err := u.validateDomainAndGetUser(user)
	if err != nil {
		return err
	}
	userDN, err := getUserDN(userInfo)
	if err != nil {
		return err
	}
	return u.ADUserAddGroup2(domain, userDN, group)

}

func (u *UseCase) ADUserAddGroup2(domain string, userDN string, groupDN string) error {
	if err := u.ad.AddUserGroup(domain, userDN, groupDN); err != nil {
		if isEntryExistsError(err) {
			return nil
		}
		return err
	}
	return nil

}

// ADUserDelGroup removes a user from a group by username and group name
func (u *UseCase) ADUserDelGroup(user string, group string) error {
	domain, userInfo, err := u.validateDomainAndGetUser(user)
	if err != nil {
		return err
	}

	userDN, err := getUserDN(userInfo)
	if err != nil {
		return err
	}

	return u.ADUserDelGroup2(domain, userDN, group)
}

// ADUserDelGroup2 removes a user from a group by domain, user DN, and group DN
// It silently succeeds if the entry does not exist
func (u *UseCase) ADUserDelGroup2(domain string, userDN string, groupDN string) error {
	if err := u.ad.DelUserGroup(domain, userDN, groupDN); err != nil {
		if isEntryNotExistsError(err) {
			return nil
		}
		return err
	}
	return nil
}

// removeUserFromAllInternetGroups removes user from all internet groups
func (u *UseCase) removeUserFromAllInternetGroups(domain, userDN string, internetGroups entity.ADInternetGroupsDN) error {
	groupsToRemove := []string{
		internetGroups.WhiteList,
		internetGroups.Full,
		internetGroups.Tech,
	}

	for _, groupDN := range groupsToRemove {
		if groupDN != "" {
			if err := u.ADUserDelGroup2(domain, userDN, groupDN); err != nil {
				return err
			}
		}
	}
	return nil
}

// getTargetGroupDN returns the target group DN based on group name
func getTargetGroupDN(groupName string, internetGroups entity.ADInternetGroupsDN) (string, error) {
	switch groupName {
	case groupNameWhitelist:
		return internetGroups.WhiteList, nil
	case groupNameFull:
		return internetGroups.Full, nil
	case groupNameTech:
		return internetGroups.Tech, nil
	case "":
		return "", nil
	default:
		return "", fmt.Errorf(errInvalidGroup, groupName)
	}
}

// getCurrentUserInternetGroup определяет текущую группу интернета пользователя
func (u *UseCase) getCurrentUserInternetGroup(user, domain string, userInfo map[string]any) string {
	internetGroupsConf := u.ad.GetDomainInternetGroups(domain)
	userGroups, ok := userInfo[fieldMemberOf].([]string)
	if !ok {
		return ""
	}

	if AnyOfFirstInSecond(internetGroupsConf.WhiteList, userGroups) {
		return groupNameWhitelist
	}
	if AnyOfFirstInSecond(internetGroupsConf.Full, userGroups) {
		return groupNameFull
	}
	if AnyOfFirstInSecond(internetGroupsConf.Tech, userGroups) {
		return groupNameTech
	}
	return ""
}

// saveTemporaryGroupChange сохраняет информацию о временном изменении группы в Redis
func (u *UseCase) saveTemporaryGroupChange(change *entity.TemporaryGroupChange) error {
	data, err := json.Marshal(change)
	if err != nil {
		return fmt.Errorf("ошибка при сериализации временного изменения: %w", err)
	}

	// Сохраняем с TTL равным времени до истечения + небольшой запас (1 час)
	ttl := time.Until(change.ExpiresAt) + time.Hour
	if ttl < 0 {
		ttl = time.Hour
	}

	// Используем существующий метод Redis для сохранения
	if err := u.redis.AddKeyFieldValue(redisKeyTempGroupChange, change.User, data); err != nil {
		return fmt.Errorf("ошибка при сохранении временного изменения: %w", err)
	}

	return nil
}

// getTemporaryGroupChange получает информацию о временном изменении группы из Redis
func (u *UseCase) getTemporaryGroupChange(user string) (*entity.TemporaryGroupChange, error) {
	data, err := u.redis.GetKeyFieldValue(redisKeyTempGroupChange, user)
	if err != nil {
		return nil, err
	}

	var change entity.TemporaryGroupChange
	if err := json.Unmarshal([]byte(data), &change); err != nil {
		return nil, fmt.Errorf("ошибка при десериализации временного изменения: %w", err)
	}

	return &change, nil
}

// GetTemporaryGroupChange получает информацию о временном изменении группы пользователя (публичный метод)
func (u *UseCase) GetTemporaryGroupChange(user string) (*entity.TemporaryGroupChange, error) {
	return u.getTemporaryGroupChange(user)
}

// deleteTemporaryGroupChange удаляет информацию о временном изменении группы из Redis
func (u *UseCase) deleteTemporaryGroupChange(user string) error {
	return u.redis.DelKeyField(redisKeyTempGroupChange, user)
}

// DeleteTemporaryGroupChange удаляет информацию о временном изменении группы пользователя (публичный метод)
// Также восстанавливает предыдущую группу пользователя
func (u *UseCase) DeleteTemporaryGroupChange(user string) error {
	// Получаем информацию о временном изменении
	change, err := u.getTemporaryGroupChange(user)
	if err != nil {
		return fmt.Errorf("временное изменение не найдено: %w", err)
	}

	// Восстанавливаем предыдущую группу
	if err := u.restoreUserGroup(change); err != nil {
		return fmt.Errorf("ошибка при восстановлении группы: %w", err)
	}

	// Удаление уже выполнено в restoreUserGroup, но для надежности удалим еще раз
	return u.deleteTemporaryGroupChange(user)
}

// restoreUserGroup восстанавливает предыдущую группу пользователя
func (u *UseCase) restoreUserGroup(change *entity.TemporaryGroupChange) error {
	internetGroups := u.ad.GetDomainInternetGroupsDN(change.Domain)

	// Удаляем пользователя из всех групп
	if err := u.removeUserFromAllInternetGroups(change.Domain, change.UserDN, internetGroups); err != nil {
		// Игнорируем ошибки доступа, так как группа может быть уже удалена
		if !isEntryAccessError(err) {
			return fmt.Errorf("ошибка при удалении из групп: %w", err)
		}
	}

	// Если была предыдущая группа, добавляем пользователя обратно
	if change.PreviousGroup != "" {
		targetGroupDN, err := getTargetGroupDN(change.PreviousGroup, internetGroups)
		if err != nil {
			return fmt.Errorf("ошибка при определении предыдущей группы: %w", err)
		}

		if targetGroupDN != "" {
			if err := u.ADUserAddGroup2(change.Domain, change.UserDN, targetGroupDN); err != nil {
				return fmt.Errorf("ошибка при восстановлении группы: %w", err)
			}
		}
	}

	// Обновляем кеш
	userInfo, err := u.ad.GetUser(change.Domain, change.User)
	if err != nil {
		// Если пользователь не найден, просто удаляем временное изменение
		u.deleteTemporaryGroupChange(change.User)
		return nil
	}

	if err := u.updateUserInternetGroupsInCache(change.User, change.Domain, userInfo[fieldMemberOf]); err != nil {
		// Логируем ошибку, но не возвращаем её, так как основная операция выполнена
		u.log.Warnf("Ошибка при обновлении кеша после восстановления группы для %s: %v", change.User, err)
	}

	// Удаляем временное изменение после успешного восстановления
	return u.deleteTemporaryGroupChange(change.User)
}

// RestoreExpiredTemporaryGroups проверяет и восстанавливает группы с истекшим временем
func (u *UseCase) RestoreExpiredTemporaryGroups() error {
	// Получаем все ключи временных изменений
	allChanges, err := u.redis.GetKeyFieldAll(redisKeyTempGroupChange)
	if err != nil {
		// Если нет данных, это не ошибка
		return nil
	}

	now := time.Now()
	var errors []error

	for user, data := range allChanges {
		var change entity.TemporaryGroupChange
		if err := json.Unmarshal([]byte(data), &change); err != nil {
			u.log.Warnf("Ошибка при десериализации временного изменения для %s: %v", user, err)
			continue
		}

		// Проверяем, истекло ли время
		if now.After(change.ExpiresAt) {
			u.log.Infof("Восстановление группы для пользователя %s (истекло в %v)", user, change.ExpiresAt)
			if err := u.restoreUserGroup(&change); err != nil {
				u.log.Errorf("Ошибка при восстановлении группы для %s: %v", user, err)
				errors = append(errors, fmt.Errorf("пользователь %s: %w", user, err))
			}
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("ошибки при восстановлении групп: %v", errors)
	}

	return nil
}

// calculateExpirationTime вычисляет время истечения для временного изменения
// days: количество суток (1 = завтра в 8:00, 2 = послезавтра в 8:00 и т.д.)
func calculateExpirationTime(days int) (time.Time, error) {
	loc, err := time.LoadLocation("Europe/Minsk")
	if err != nil {
		return time.Time{}, fmt.Errorf("ошибка загрузки временной зоны: %w", err)
	}

	now := time.Now().In(loc)
	// Возвращаем указанное количество дней от сегодня в 8:00 утра
	// days=1: завтра в 8:00, days=2: послезавтра в 8:00 и т.д.
	return time.Date(now.Year(), now.Month(), now.Day(), 8, 0, 0, 0, loc).AddDate(0, 0, days), nil
}

// SwitchUserGroupInternet switches a user's internet group membership
// isTemporary: если true, изменение будет временным и автоматически вернется обратно
// days: количество суток (1 = завтра в 8:00, 2 = послезавтра в 8:00 и т.д., используется только если isTemporary = true)
func (u *UseCase) SwitchUserGroupInternet(tuser string, user string, groupName string, isTemporary bool, days int) error {
	// Validate domain and get user info
	domain, userInfo, err := u.validateDomainAndGetUser(user)
	if err != nil {
		return err
	}
	// Check permissions
	if err := u.UserInDomainGroup2(tuser, permissionGroupDomainAdmins, domain); err != nil {
		return u.Error(errNoPermission)
	}

	userDN, err := getUserDN(userInfo)
	if err != nil {
		return u.Error(err.Error())
	}

	// Определяем текущую группу пользователя (для временного режима)
	var previousGroup string
	if isTemporary {
		previousGroup = u.getCurrentUserInternetGroup(user, domain, userInfo)
	}

	internetGroups := u.ad.GetDomainInternetGroupsDN(domain)
	// Remove user from all internet groups
	if err := u.removeUserFromAllInternetGroups(domain, userDN, internetGroups); err != nil {
		//INSUFF_ACCESS_RIGHTS
		if isEntryAccessError(err) {
			return u.Error("Отсутствуют права доступа.")
		}
		return err
	}

	// Получаем DN целевой группы по имени
	targetGroupDN, err := getTargetGroupDN(groupName, internetGroups)
	if err != nil {
		return err
	}
	// If no target group specified, we're done (user removed from all groups)
	if targetGroupDN != "" {
		// Add user to target group
		if err := u.ADUserAddGroup2(domain, userDN, targetGroupDN); err != nil {
			return fmt.Errorf(errFailedToAddToGroup, groupName, err)
		}
	}

	usercn, _ := getUsercn(userInfo)
	domainAdminsGLPIId := u.ad.GetDomainAdminsGLPI(domain)
	_, channelId, _, _ := u.glpi.GetGroupMattermostChannel(domainAdminsGLPIId)

	// Если это временное изменение, сохраняем информацию для автоматического возврата
	if isTemporary {
		expiresAt, err := calculateExpirationTime(days)
		if err != nil {
			u.log.Warnf("Ошибка при расчете времени истечения для %s: %v", user, err)
			// Используем завтра в 8:00 по умолчанию
			loc, _ := time.LoadLocation("Europe/Minsk")
			now := time.Now().In(loc)
			expiresAt = time.Date(now.Year(), now.Month(), now.Day(), 8, 0, 0, 0, loc).AddDate(0, 0, days)

		}

		change := &entity.TemporaryGroupChange{
			User:          user,
			Domain:        domain,
			UserDN:        userDN,
			PreviousGroup: previousGroup,
			NewGroup:      groupName,
			CreatedAt:     time.Now(),
			ExpiresAt:     expiresAt,
			ChangedBy:     tuser,
		}

		if err := u.saveTemporaryGroupChange(change); err != nil {
			u.log.Warnf("Не удалось сохранить информацию о временном изменении для %s: %v", user, err)
			// Не возвращаем ошибку, так как основная операция выполнена успешно
		}
		message := fmt.Sprintf(`У пользователя %s (%s) с помощью UserInfo **временно** изменилась /n
		текущая группа интернет с %s на %s /n
        группа %s вернётся у пользователя %s /n
		*Группу изменил %s*
		`, usercn, user, previousGroup, groupName, previousGroup, expiresAt, tuser)
		if len(channelId) > 0 {
			u.SendPostSimple(channelId, message)
		}
	} else {
		message := fmt.Sprintf(`У пользователя %s (%s) с помощью UserInfo изменилась текущая группа интернет на %s`, usercn, user, groupName)
		if len(channelId) > 0 {
			u.SendPostSimple(channelId, message)
		}
	}

	// Refresh user info from AD
	userInfo, err = u.ad.GetUser(domain, user)
	if err != nil {
		return u.Error(errUserNotFound)
	}

	// Update cache
	return u.updateUserInternetGroupsInCache(user, domain, userInfo[fieldMemberOf])
}

// updateUserInternetGroupsInCache updates user internet group flags in Redis cache
func updateUserInternetGroups(cachedUser map[string]any, memberOf any, internetGroupsConf entity.ADInternetGroups) map[string]any {
	// Update memberOf from fresh AD data
	cachedUser[fieldMemberOf] = memberOf
	// Get internet groups configuration
	// Extract user groups
	userGroups, ok := cachedUser[fieldMemberOf].([]string)
	if !ok {
		userGroups = []string{}
	}
	// Remove old internet group flags
	delete(cachedUser, fieldInternetWL)
	delete(cachedUser, fieldInternet)
	delete(cachedUser, fieldInternetTech)

	// Set new internet group flags based on membership
	if AnyOfFirstInSecond(internetGroupsConf.WhiteList, userGroups) {
		cachedUser[fieldInternetWL] = true
	}
	if AnyOfFirstInSecond(internetGroupsConf.Full, userGroups) {
		cachedUser[fieldInternet] = true
	}
	if AnyOfFirstInSecond(internetGroupsConf.Tech, userGroups) {
		cachedUser[fieldInternetTech] = true
	}
	return cachedUser
}

// updateUserInternetGroupsInCache updates user internet group flags in Redis cache
func (u *UseCase) updateUserInternetGroupsInCache(user, domain string, memberOf any) error {
	userJSON, err := u.redis.GetKeyFieldValue(redisKeyAllUsers, user)
	if err != nil {
		return u.Error(errUserNotInSystem)
	}

	var cachedUser map[string]any
	if err := unmarshalString(userJSON, &cachedUser); err != nil {
		return fmt.Errorf("ошибка при обновлении кеша, ошибка парсинга JSON: %w", err)
	}

	// Get internet groups configuration
	internetGroupsConf := u.ad.GetDomainInternetGroups(domain)

	cachedUser = updateUserInternetGroups(cachedUser, memberOf, internetGroupsConf)

	// Marshal and save to Redis
	jsonUser, err := json.Marshal(cachedUser)
	if err != nil {
		return fmt.Errorf("ошибка при обновлении кеша, ошибка парсинга JSON: %w", err)
	}

	upn, ok := cachedUser[fieldUserPrincipalName].(string)
	if ok && upn != "" {
		if err := u.redis.AddKeyFieldValue(redisKeyAllUsers, upn, jsonUser); err != nil {
			return fmt.Errorf("ошибка при обновлении кеша: %w", err)
		}
	}

	// Update all domain cash AD data
	oneDomain, err := u.redis.GetKeyFieldValue("ad", domain)
	var users []map[string]interface{}
	if err := json.Unmarshal([]byte(oneDomain), &users); err != nil {
		return fmt.Errorf("ошибка при обновлении кеша, ошибка парсинга JSON: %w", err)
	}
	for i, newuser := range users {
		if userName, ok := newuser["userPrincipalName"].(string); ok && userName == user {
			users[i] = updateUserInternetGroups(newuser, memberOf, internetGroupsConf)
			break
		}
	}
	updatedBytes, err := json.Marshal(users)
	if err := u.redis.AddKeyFieldValue("ad", domain, updatedBytes); err != nil { // "" для полного хеша
		return fmt.Errorf("ошибка при обновлении кеша Redis: %w", err)
	}

	return nil
}
