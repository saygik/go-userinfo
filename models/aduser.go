package models

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/saygik/go-userinfo/ad"
	"github.com/saygik/go-userinfo/db"
)

//UserModel ...
type ADUserModel struct{}

var ctx = context.Background()

//AllDomains...
func (m ADUserModel) AllDomains() map[string]bool {
	return ad.Domains
}

//type User struct {
//	UserPrincipalName string `db:"userPrincipalName" json:"userPrincipalName"`
//	Dn string `db:"dn" json:"dn"`
//	Cn string `db:"cn" json:"cn"`
//	Company string `db:"company" json:"company"`
//	Department string `db:"department" json:"department"`
//	Title string `db:"title" json:"title"`
//	TelephoneNumber string `db:"telephoneNumber" json:"telephoneNumber"`
//	OtherTelephone string `db:"otherTelephone" json:"otherTelephone"`
//	Mobile string `db:"mobile" json:"mobile"`
//	Mail string `db:"mail" json:"mail"`
//	Pager string `db:"pager" json:"pager"`
//	Sip string `db:"msRTCSIP-PrimaryUserAddress" json:"msRTCSIP-PrimaryUserAddress"`
//	Url string `db:"url" json:"url"`
//	MemberOf string `db:"memberOf" json:"memberOf"`
//}
//"userPrincipalName", "dn", "cn", "company", "department", "title", "telephoneNumber",	"otherTelephone", "mobile", "mail", "pager", "msRTCSIP-PrimaryUserAddress", "url","memberOf"
//All ...
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
	ips, err := UserModel{}.All(domain)
	presences, err := SkypeModel{}.AllPresences()
	if err == nil {
		if len(ips) > 0 || len(presences) > 0 {
			for _, user := range users {
				if isStringInArray("Пользователи интернета", user["memberOf"]) {
					user["internet"] = true
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
							user["presencetime"] = presence.Lastpubtime
						}
					}
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
	users, err := ad.GetAD(domain).GetGroupUsers("CN=Пользователи интернета,OU=_Groups,DC=brnv,DC=rw")
	if err != nil || len(users) < 1 {
		return nil, errors.New("Users not found")
	}

	return users, err
}
