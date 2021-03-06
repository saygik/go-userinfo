package ad

import (
	adClient "github.com/saygik/go-ad-client"
	"time"
)

var adClients = map[string]*adClient.ADClient{}

type ADConfig struct {
	Key          string        `json:"key"`
	Name         string        `json:"name"`
	Base         string        `json:"base"`
	Dc           string        `json:"dc"`
	Filter       string        `json:"filter"`
	BindDN       string        `json:"bindDN"`
	BindPassword string        `json:"bindPassword"`
	Time         time.Duration `json:"time"`
}
type Config struct {
	ADS []ADConfig
}

//One minute ttl 60000000000
const AllUsersTTL time.Duration = 60000000000

var Domains map[string]bool

func Init(adconfig Config) {
	Domains = make(map[string]bool)
	for _, oneADConfig := range adconfig.ADS {
		Domains[oneADConfig.Key] = true
		//		Domains=append(Domains,oneADConfig.Key)
		oneADClient := NewAddConnection(oneADConfig)
		defer oneADClient.Close()
		adClients[oneADConfig.Key] = oneADClient
	}
	users, err := GetAD("brnv.rw").GetAllUsers()
	if err != nil || len(users) < 1 {
		return
	}

}

//"(userPrincipalName=%s)"
func NewAddConnection(config ADConfig) *adClient.ADClient {
	client := &adClient.ADClient{
		Base:         config.Base,
		Host:         config.Dc,
		Port:         389,
		UseSSL:       false,
		BindDN:       config.BindDN,
		BindPassword: config.BindPassword,
		UserFilter:   config.Filter,
		GroupFilter:  "(memberUid=%s)",
		Attributes: []string{"userPrincipalName", "dn", "cn", "company", "department", "title", "telephoneNumber",
			"otherTelephone", "mobile", "mail", "pager", "msRTCSIP-PrimaryUserAddress", "url"},
	}
	return client
}

func GetAD(adName string) *adClient.ADClient {
	cli := adClients[adName]
	return cli
}

func Close() {
	//	client.Close()
}
