package models

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/saygik/go-userinfo/ad"
	"github.com/saygik/go-userinfo/db"
	"github.com/saygik/go-userinfo/forms"
)

// UserModel ...
type ADUserModel struct{}

var ctx = context.Background()

// AllDomains...
func (m ADUserModel) AllDomains() []ad.ADArray {
	return ad.DomainsArray
}

//	type User struct {
//		UserPrincipalName string `db:"userPrincipalName" json:"userPrincipalName"`
//		Dn string `db:"dn" json:"dn"`
//		Cn string `db:"cn" json:"cn"`
//		Company string `db:"company" json:"company"`
//		Department string `db:"department" json:"department"`
//		Title string `db:"title" json:"title"`
//		TelephoneNumber string `db:"telephoneNumber" json:"telephoneNumber"`
//		OtherTelephone string `db:"otherTelephone" json:"otherTelephone"`
//		Mobile string `db:"mobile" json:"mobile"`
//		Mail string `db:"mail" json:"mail"`
//		Pager string `db:"pager" json:"pager"`
//		Sip string `db:"msRTCSIP-PrimaryUserAddress" json:"msRTCSIP-PrimaryUserAddress"`
//		Url string `db:"url" json:"url"`
//		MemberOf string `db:"memberOf" json:"memberOf"`
//	}
//
// "userPrincipalName", "dn", "cn", "company", "department", "title", "telephoneNumber",	"otherTelephone", "mobile", "mail", "pager", "msRTCSIP-PrimaryUserAddress", "url","memberOf"

func (m ADUserModel) ClearAllDomainsUsers() {
	redisClient := db.GetRedis()
	redisClient.Del(ctx, "ad")
	redisClient.Del(ctx, "adc")
}

func (m ADUserModel) GetAllDomainsUsers() {
	allADs := ad.GetAllADClients()
	redisClient := db.GetRedis()

	for domain, oneAD := range allADs {
		users, err := oneAD.GetAllUsers()
		comps, _ := oneAD.GetAllComputers()

		if err == nil || len(users) > 0 {
			//break // break here
			println("Get from ad to redis from " + domain)
			ips, _ := UserIPModel{}.All(domain)
			avatars, _ := UserIPModel{}.AllAvatars(domain)
			presences, _ := SkypeModel{}.AllPresences()
			for _, user := range users {
				user["domain"] = domain
				if IsStringInArray("Пользователи интернета", user["memberOf"]) {
					user["internet"] = true
				}
				if IsStringInArray("Пользователи интернета Белый список", user["memberOf"]) {
					user["internetwl"] = true
				}
				if len(ips) > 0 {
					for _, ip := range ips {
						if user["userPrincipalName"] == ip.Login {
							user["ip"] = ip.Ip
							user["computer"] = ip.Computer
						}
					}
				}
				if len(avatars) > 0 {
					for _, avatar := range avatars {
						if user["userPrincipalName"] == avatar.Name {
							user["avatar"] = avatar.Avatar
						}
					}
				}

				// if _, ok := user["avatar"]; ok {
				// 	user["avatar"] = "avatar1"
				// }
				if len(presences) > 0 {
					for _, presence := range presences {
						if user["userPrincipalName"] == presence.Userathost {
							user["presence"] = presence.Presence
							user["lastpubtime"] = presence.Lastpubtime

						}

					}
				}
				jsonUser, _ := json.Marshal(user)
				redisClient.HSet(ctx, "allusers", user["userPrincipalName"], jsonUser).Err()
			}

		}
		sort.Slice(users, func(i, j int) bool {
			return fmt.Sprintf("%v", users[i]["cn"]) < fmt.Sprintf("%v", users[j]["cn"])
		})
		jsonUsers, _ := json.Marshal(users)
		jsonComps, _ := json.Marshal(comps)
		err1 := redisClient.HSet(ctx, "ad", domain, jsonUsers).Err()
		if err1 != nil {
			fmt.Println("key does not exists")
			return
		}
		redisClient.HSet(ctx, "adc", domain, jsonComps).Err()

	}
}

func (m ADUserModel) AllAdCounts() (users int, computers int, error1 error) {
	var r []map[string]interface{}

	redisClient := db.GetRedis()
	if redisClient == nil {
		return 0, 0, errors.New("Redis not found")
	}
	redisADUsers, err := redisClient.HGetAll(ctx, "ad").Result()
	if err != nil {
		return 0, 0, err
	}
	redisADComputers, err := redisClient.HGetAll(ctx, "adc").Result()
	if err != nil {
		return 0, 0, err
	}
	users = 0
	computers = 0
	for _, oneDomain := range redisADUsers {
		json.Unmarshal([]byte(oneDomain), &r)
		users = users + len(r)
	}
	for _, oneDomain := range redisADComputers {
		json.Unmarshal([]byte(oneDomain), &r)
		computers = computers + len(r)
	}

	return users, computers, nil

}

// All Computers...
func (m ADUserModel) AllAdComputers(user string) ([]map[string]interface{}, error) {

	var res []map[string]interface{}
	redisClient := db.GetRedis()
	if redisClient == nil {
		return nil, errors.New("Redis not found")
	}
	redisADUsers, err := redisClient.HGetAll(ctx, "adc").Result()
	if err != nil {
		return nil, err
	}
	for domainName, oneDomain := range redisADUsers {
		access := GetAccessToResource(domainName, user)
		if access == -1 {
			continue
		}
		var r []map[string]interface{}
		json.Unmarshal([]byte(oneDomain), &r)

		res = append(res, r...)
	}
	return res, nil
}

// All Ad Users Short info...
func (m ADUserModel) AllAdUsersShort() ([]map[string]interface{}, error) {
	redisClient := db.GetRedis()
	var ctx = context.Background()
	var users []map[string]interface{}

	redisADUsers, err := redisClient.HGetAll(ctx, "allusers").Result()
	if err != nil {
		return []map[string]interface{}{}, err
	}
	for _, value := range redisADUsers {
		var user map[string]interface{}
		json.Unmarshal([]byte(value), &user)
		delete(user, "ip")
		delete(user, "pwdLastSet")
		delete(user, "proxyAddresses")
		delete(user, "passwordDontExpire")
		delete(user, "passwordCantChange")
		delete(user, "distinguishedName")
		delete(user, "userAccountControl")
		delete(user, "memberOf")
		delete(user, "employeeNumber")
		delete(user, "presence")
		delete(user, "url")
		delete(user, "otherTelephone")
		user["findedInAD"] = true
		user["name"] = user["userPrincipalName"]
		users = append(users, user)
	}
	//	json.Unmarshal([]byte(redisADUsers), &users)
	return users, nil
}

// All ...
func (m ADUserModel) AllAd(user string) ([]map[string]interface{}, error) {
	domain := strings.Split(fmt.Sprintf("%s", user), "@")[1]
	redisClient := db.GetRedis()
	if redisClient == nil {
		return nil, errors.New("Redis not found")
	}

	//redisADUsers, err := redisClient.Get(ctx, "brnv.rw"+"-ad").Result()
	redisADUsers, err := redisClient.HGetAll(ctx, "ad").Result()
	if err != nil {
		return nil, err
	}

	var res []map[string]interface{}
	//res := make([]map[string]interface{}, 3000)
	for domainName, oneDomain := range redisADUsers {
		access := GetAccessToResource(domainName, user)
		if access == -1 && domain == domainName {
			access = 0
		}
		if access == -1 {
			continue
		}
		var r []map[string]interface{}
		json.Unmarshal([]byte(oneDomain), &r)
		// isUserAccessToDomain := IsStringInArray(domainName, userRoles) || IsStringInArray("fullAdmin", userRoles)
		// domainTechnical := IsStringInArray("domainTechnical", userRoles) || IsStringInArray("domainAdmin", userRoles)
		//accessToTechnicalInfo := (isUserAccessToDomain && domainTechnical) || IsStringInArray("fullAdmin", userRoles) || access
		accessToTechnicalInfo := access == 1
		for _, user := range r {
			delete(user, "employeeNumber")
			if !accessToTechnicalInfo {
				delete(user, "ip")
				delete(user, "pwdLastSet")
				delete(user, "proxyAddresses")
				delete(user, "passwordDontExpire")
				delete(user, "passwordCantChange")
				delete(user, "distinguishedName")
				delete(user, "userAccountControl")
				delete(user, "memberOf")

				user["restricted"] = true
			}

			//user["ip"] = "-"
		}
		res = append(res, r...)
	}
	// if err == nil && redisADUsers != "" {
	// 	if redisADUsers != "" && redisADUsers != "null" {
	// 		var r []map[string]interface{}
	// 		json.Unmarshal([]byte(fmt.Sprintf("%v", redisADUsers)), &r)
	// 		println("Get from redis")
	// 		return r, nil
	// 	}
	// }
	return res, nil
}

// All ...
func (m ADUserModel) All(domain string) ([]map[string]interface{}, error) {
	//ad.GetDomainUsers(domain)
	//return nil,nil

	redisClient := db.GetRedis()
	if redisClient == nil {
		return nil, errors.New("Redis not found")
	}
	redisADUsers, err := redisClient.Get(ctx, domain+"-ad").Result()

	if err == nil && redisADUsers != "" {
		if redisADUsers != "" && redisADUsers != "null" {
			var r []map[string]interface{}
			json.Unmarshal([]byte(redisADUsers), &r)
			println("Get from redis")
			return r, nil
		}
	}
	return nil, nil
	currentADclient := ad.GetAD(domain)
	if currentADclient == nil {
		return nil, errors.New("This domain is not served by the system")
	}
	users, err := ad.GetAD(domain).GetAllUsers()
	if err != nil || len(users) < 1 {
		return nil, errors.New("Users not found")
	}
	//t := time.Now()
	//fmt.Println(t.Format("20060102150405"))
	ips, err := UserIPModel{}.All(domain)
	avatars, _ := UserIPModel{}.AllAvatars(domain)

	presences, err := SkypeModel{}.AllPresences()

	for _, user := range users {
		user["domain"] = domain
		if IsStringInArray("Пользователи интернета", user["memberOf"]) {
			user["internet"] = true
		}
		if IsStringInArray("Пользователи интернета Белый список", user["memberOf"]) {
			user["internetwl"] = true
		}
		delete(user, "memberOf")
		if len(ips) > 0 {
			for _, ip := range ips {
				if user["userPrincipalName"] == ip.Login {
					user["ip"] = ip.Ip
				}
			}
		}
		if len(avatars) > 0 {
			for _, avatar := range avatars {
				if user["userPrincipalName"] == avatar.Name {
					user["avatar"] = avatar.Avatar
				}
			}
		}
		if len(presences) > 0 {
			for _, presence := range presences {
				if user["userPrincipalName"] == presence.Userathost {
					user["presence"] = presence.Presence
					user["lastpubtime"] = presence.Lastpubtime

				}

			}
		}
	}
	//t = time.Now()
	//fmt.Println(t.Format("20060102150405"))

	jsonUsers, _ := json.Marshal(users)
	redisClient.Set(ctx, domain+"-ad", jsonUsers, ad.AllUsersTTL)
	return users, err
}

func (m ADUserModel) GroupUsers(domain string, group string) ([]map[string]interface{}, error) {

	currentADclient := ad.GetAD(domain)
	if currentADclient == nil {
		return nil, errors.New("This domain is not served by the system")
	}
	users, err := ad.GetAD(domain).GetGroupUsers(group)
	if err != nil || len(users) < 1 {
		return nil, errors.New("Users not found")
	}

	return users, err
}
func (m ADUserModel) Authenticate(form forms.LoginForm) (bool, map[string]string, error) {
	domain := strings.Split(fmt.Sprintf("%s", form.Email), "@")[1]
	catalog := ad.GetAD(domain)
	if catalog == nil {
		return false, nil, errors.New("there is no such domain")
	}
	return catalog.Authenticate(form.Email, form.Password)
}
func (m ADUserModel) GetCurrentUser(username string) (map[string]interface{}, error) {
	domain := strings.Split(fmt.Sprintf("%s", username), "@")[1]
	catalog := ad.GetAD(domain)
	if catalog == nil {
		return nil, errors.New("there is no such domain")
	}
	user, err := ad.GetAD(domain).GetUserInfo(username)
	if err != nil {
		return nil, err
	}
	var userIPModel = new(UserIPModel)
	role := IdName{Id: 5, Name: "Пользователь"}
	roles, err := userIPModel.GetUserRoles(username)
	if err == nil && len(roles) > 0 {
		role = roles[0]
	}

	groups, err := userIPModel.GetUserGroups(username)
	if err == nil && len(groups) > 0 {
		user["groups"] = groups
	} else {
		user["groups"] = []IdName{}
	}
	avatar, err := userIPModel.GetUserAvatar(username)
	if err == nil {
		user["avatar"] = avatar
	}

	user["role"] = role
	user["domain"] = domain

	delete(user, "memberOf")
	return user, nil
}

func (m ADUserModel) GetOneUserPropertys(username string, techUser string) (map[string]interface{}, error) {

	domain := GetDomainFromUserName(username)
	domainName := GetDomainFromUserName(techUser)

	access := GetAccessToResource(domain, techUser)
	if access == -1 && domain == domainName {
		access = 0
	}
	if access == -1 {
		return nil, errors.New("У вас недостаточно прав на просмотр данных пользователя")
	}

	accessToTechnicalInfo := access == 1

	catalog := ad.GetAD(domain)
	if catalog == nil {
		return nil, errors.New("нет такого домена")
	}
	redisClient := db.GetRedis()
	if redisClient == nil {
		return nil, errors.New("Redis не найден на сервере")
	}
	redisADUsers, err := redisClient.HGetAll(ctx, "ad").Result()
	if err != nil {
		return nil, err
	}
	var userIPModel = new(UserIPModel)
	ADUsers := redisADUsers[domain]

	if err == nil && ADUsers != "" {
		if ADUsers != "" && ADUsers != "null" {
			var r []map[string]interface{}
			json.Unmarshal([]byte(ADUsers), &r)
			user := FindUserInRedisArray(r, username)
			if user != nil {
				if !accessToTechnicalInfo {
					delete(user, "ip")
					delete(user, "pwdLastSet")
					delete(user, "proxyAddresses")
					delete(user, "employeeNumber")
					delete(user, "passwordDontExpire")
					delete(user, "passwordCantChange")
					delete(user, "distinguishedName")
					delete(user, "userAccountControl")
					delete(user, "memberOf")

				}

				avatar, err := userIPModel.GetUserAvatar(username)
				if err == nil {
					user["avatar"] = avatar
				}
				return user, nil
			}
		}
	}
	return nil, errors.New("пользователь не найден")
}
