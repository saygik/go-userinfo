package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

// UserSessionInfo ...
type UserSessionInfo struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// JSONRaw ...
type JSONRaw json.RawMessage

// Value ...
func (j JSONRaw) Value() (driver.Value, error) {
	byteArr := []byte(j)
	return driver.Value(byteArr), nil
}

// Scan ...
func (j *JSONRaw) Scan(src interface{}) error {
	asBytes, ok := src.([]byte)
	if !ok {
		return error(errors.New("Scan source was not []bytes"))
	}
	err := json.Unmarshal(asBytes, &j)
	if err != nil {
		return error(errors.New("Scan could not unmarshal to []string"))
	}
	return nil
}

// MarshalJSON ...
func (j *JSONRaw) MarshalJSON() ([]byte, error) {
	return *j, nil
}

// UnmarshalJSON ...
func (j *JSONRaw) UnmarshalJSON(data []byte) error {
	if j == nil {
		return errors.New("json.RawMessage: UnmarshalJSON on nil pointer")
	}
	*j = append((*j)[0:0], data...)
	return nil
}

// DataList ....
type DataList struct {
	Data JSONRaw `db:"data" json:"data"`
	Meta JSONRaw `db:"meta" json:"meta"`
}

func IsStringInArray(str string, arr interface{}) bool {
	if arr == nil {
		return false
	}
	for _, b := range arr.([]string) {
		if b == str {
			return true
		}
	}
	return false
}
func FindUserInRedisArray(users []map[string]interface{}, userToFind string) map[string]interface{} {
	for _, user := range users {
		if user["userPrincipalName"] == userToFind {
			return user
		}
	}
	return nil
}

func IsStringInArrayIdName(str string, arr []IdName) bool {
	if arr == nil {
		return false
	}
	for _, b := range arr {
		if b.Name == str {
			return true
		}
	}
	return false
}
func GetAccessToResource(resource string, user string) (res int) {
	domain := strings.Split(fmt.Sprintf("%s", user), "@")[1]
	var userIPModel = new(UserIPModel)
	res = -1
	userRoles, err := userIPModel.GetUserRoles(user)
	if err != nil {
		return -1
	}
	if IsStringInArrayIdName("Администратор системы", userRoles) {
		return 1
	}
	if IsStringInArrayIdName("Администратор", userRoles) {
		return 1
	}
	if IsStringInArrayIdName("Технический специалист", userRoles) && domain == resource {
		return 1
	}
	accessRole, _ := userIPModel.GetUserResourceAccess(resource, user)

	return accessRole
}

func GetDomainFromUserName(s string) string {
	return strings.Split(fmt.Sprintf("%s", s), "@")[1]
}
