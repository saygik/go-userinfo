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

//All ...
func (m ADUserModel) All(domain string) ([]map[string]string, error) {
	redisClient := db.GetRedis()
	if redisClient == nil {
		return nil, errors.New("Redis not found")
	}
	redisADUsers, err := redisClient.Get(ctx, domain+"-ad").Result()

	if err == nil && redisADUsers != "" {
		if redisADUsers != "" && redisADUsers != "null" {
			var r []map[string]string
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
	if err == nil {
		if len(ips) > 0 {
			for _, user := range users {
				for _, ip := range ips {
					if user["userPrincipalName"] == ip.Login {
						user["ip"] = ip.Ip
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
