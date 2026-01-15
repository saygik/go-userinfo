package usecase

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/saygik/go-userinfo/internal/entity"
)

const (
	// Error messages
	errDomainNotFound     = "домен не найден"
	errUserNotFound       = "пользователь не существует"
	errUserDNNotDefined   = "не орпеделен dn пользователя"
	errNoPermission       = "У Вас нет прав на эту операцию"
	errUserNotInSystem    = "нет таккого пользователя в системе"
	errCashAllNotFound    = "нет общего кеша для обновления"
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
	redisKeyAllUsers = "allusers"

	// User fields
	fieldDistinguishedName = "distinguishedName"
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
	// err = u.ad.AddUserGroup(domain, userDN, group)
	// if err != nil {
	// 	//		"ENTRY_EXISTS"
	// 	if strings.Contains(err.Error(), "ENTRY_EXISTS") || strings.Contains(err.Error(), "already exists") {
	// 		return nil
	// 	}
	// 	return u.Error(err.Error())
	// }
	// return nil
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

func (u *UseCase) SwitchUserGroupInternet(tuser string, user string, groupName string) error {

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
