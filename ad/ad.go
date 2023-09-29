package ad

import (
	"time"

	adClient "github.com/saygik/go-ad-client"
)

var adClients = map[string]*adClient.ADClient{}

type ADConfig struct {
	Key          string        `json:"key"`
	Name         string        `json:"name"`
	Base         string        `json:"base"`
	Dc           string        `json:"dc"`
	GroupFilter  string        `json:"group-filter"`
	Filter       string        `json:"filter"`
	BindDN       string        `json:"bindDN"`
	BindPassword string        `json:"bindPassword"`
	Time         time.Duration `json:"time"`
}
type ADArray struct {
	Name  string `json:"name"`
	Title string `json:"title"`
}
type Config struct {
	ADS []ADConfig
}

type User struct {
	UserPrincipalName []string `db:"userPrincipalName" json:"userPrincipalName"`
	Dn                string   `db:"dn" json:"dn"`
	Cn                string   `db:"cn" json:"cn"`
}

// One minute ttl 60000000000
const AllUsersTTL time.Duration = 300000000000

var Domains map[string]bool

var DomainsArray []ADArray

func Init(adconfig Config) {
	Domains = make(map[string]bool)
	for _, oneADConfig := range adconfig.ADS {
		domain := ADArray{Name: oneADConfig.Key, Title: oneADConfig.Name}
		DomainsArray = append(DomainsArray, domain)
		Domains[oneADConfig.Key] = true
		//		Domains=append(Domains,oneADConfig.Key)
		oneADClient := NewAddConnection(oneADConfig)
		oneADClient.SkipTLS = true
		defer oneADClient.Close()
		adClients[oneADConfig.Key] = oneADClient
	}
	//	users, err := GetAD("brnv.rw").GetAllUsers()
	//	if err != nil || len(users) < 1 {
	//		return
	//	}

}

// "(userPrincipalName=%s)"
func NewAddConnection(config ADConfig) *adClient.ADClient {
	client := &adClient.ADClient{
		Base:         config.Base,
		Host:         config.Dc,
		Port:         389,
		UseSSL:       false,
		BindDN:       config.BindDN,
		BindPassword: config.BindPassword,
		UserFilter:   config.Filter,
		GroupFilter:  config.GroupFilter,
		Attributes: []string{"userPrincipalName", "dn", "cn", "company", "department", "title", "telephoneNumber",
			"otherTelephone", "mobile", "mail", "pager", "msRTCSIP-PrimaryUserAddress", "url", "memberOf", "displayName",
			"description", "userPrincipalName", "employeeNumber", "pwdLastSet", "proxyAddresses", "userAccountControl", "distinguishedName"},
	}
	return client
}

func GetAD(adName string) *adClient.ADClient {
	cli := adClients[adName]
	return cli
}

func GetAllADClients() map[string]*adClient.ADClient {
	return adClients
}

func Close() {
	//	client.Close()
}
