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
}
func (m ADUserModel) GetAllDomainsUsers() {
	allADs := ad.GetAllADClients()
	redisClient := db.GetRedis()
	for domain, oneAD := range allADs {
		users, err := oneAD.GetAllUsers()
		if err == nil || len(users) > 0 {
			//break // break here
			println("Get from ad to redis from " + domain)
			ips, _ := UserIPModel{}.All(domain)
			presences, _ := SkypeModel{}.AllPresences()
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
				if len(presences) > 0 {
					for _, presence := range presences {
						if user["userPrincipalName"] == presence.Userathost {
							user["presence"] = presence.Presence
							user["lastpubtime"] = presence.Lastpubtime

						}

					}
				}
			}

		}
		sort.Slice(users, func(i, j int) bool {
			return fmt.Sprintf("%v", users[i]["cn"]) < fmt.Sprintf("%v", users[j]["cn"])
		})
		jsonUsers, _ := json.Marshal(users)
		err1 := redisClient.HSet(ctx, "ad", domain, jsonUsers).Err()
		if err1 != nil {
			fmt.Println("key does not exists")
			return
		}
	}
}

// All ...
func (m ADUserModel) AllAd(userRoles []string) ([]map[string]interface{}, error) {

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
		var r []map[string]interface{}
		json.Unmarshal([]byte(oneDomain), &r)
		isUserAccessToDomain := IsStringInArray(domainName, userRoles) || IsStringInArray("fullAdmin", userRoles)
		domainTechnical := IsStringInArray("domainTechnical", userRoles) || IsStringInArray("domainAdmin", userRoles)
		accessToTechnicalInfo := (isUserAccessToDomain && domainTechnical) || IsStringInArray("fullAdmin", userRoles)
		for _, user := range r {
			if !accessToTechnicalInfo {
				delete(user, "ip")
				delete(user, "pwdLastSet")
				delete(user, "proxyAddresses")
				delete(user, "employeeNumber")
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
func (m ADUserModel) GetOneUser(username string) (map[string]interface{}, error) {
	domain := strings.Split(fmt.Sprintf("%s", username), "@")[1]
	catalog := ad.GetAD(domain)
	if catalog == nil {
		return nil, errors.New("there is no such domain")
	}
	user, err := ad.GetAD(domain).GetUserInfo(username)
	if err != nil {
		return nil, err
	}

	user["domain"] = domain
	if IsStringInArray("adusersGlobalAdmins", user["memberOf"]) && domain == "brnv.rw" {
		user["role"] = "globaladmin"
		delete(user, "memberOf")
		return user, nil
	}
	if IsStringInArray("adusersDomainAdmins", user["memberOf"]) {
		user["role"] = "admin"
		delete(user, "memberOf")
		return user, nil
	}
	if IsStringInArray("adusersTS", user["memberOf"]) {
		user["role"] = "ts"
		delete(user, "memberOf")
		return user, nil
	}
	user["role"] = "user"
	delete(user, "memberOf")
	return user, nil
}

func (m ADUserModel) GetOneUserPropertys(username string, userRoles []string) (map[string]interface{}, error) {
	domain := strings.Split(fmt.Sprintf("%s", username), "@")[1]
	isUserAccessToDomain := IsStringInArray(domain, userRoles) || IsStringInArray("fullAdmin", userRoles)
	domainTechnical := IsStringInArray("domainTechnical", userRoles) || IsStringInArray("domainAdmin", userRoles)
	accessToTechnicalInfo := (isUserAccessToDomain && domainTechnical) || IsStringInArray("fullAdmin", userRoles)

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
				}
				return user, nil
			}
		}
	}
	return nil, errors.New("пользователь не найден")
}

func (m ADUserModel) GetAdusersRights(username string) (role string, err error) {
	role = "user"
	domain := strings.Split(fmt.Sprintf("%s", username), "@")[1]
	catalog := ad.GetAD(domain)
	if catalog == nil {
		return role, errors.New("there is no such domain")
	}
	user, err := ad.GetAD(domain).GetUserInfo(username)
	if err != nil {
		return role, err
	}

	if IsStringInArray("adusersGlobalAdmins", user["memberOf"]) && domain == "brnv.rw" {
		role = "globaladmin"
	}
	if IsStringInArray("adusersDomainAdmins", user["memberOf"]) {
		role = "admin"
	}
	if IsStringInArray("adusersTS", user["memberOf"]) {
		role = "ts"
	}
	return role, nil
}
