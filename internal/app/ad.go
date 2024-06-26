package app

import (
	adClient "github.com/saygik/go-ad-client"
	"github.com/saygik/go-userinfo/internal/config"
)

// type adArray struct {
// 	Name  string `json:"name"`
// 	Title string `json:"title"`
// }

func NewADClients(adConfigs []config.ADConfig) map[string]*adClient.ADClient {
	var adClients = map[string]*adClient.ADClient{}
	for _, oneADConfig := range adConfigs {
		oneADClient := newAddConnection(oneADConfig)
		oneADClient.SkipTLS = true
		defer oneADClient.Close()
		adClients[oneADConfig.Key] = oneADClient
	}
	return adClients
}

func newAddConnection(config config.ADConfig) *adClient.ADClient {
	client := &adClient.ADClient{
		Title:          config.Name,
		Domain:         config.Key,
		Base:           config.Base,
		Host:           config.Dc,
		Port:           389,
		UseSSL:         false,
		BindDN:         config.BindDN,
		BindPassword:   config.BindPassword,
		UserFilter:     config.Filter,
		ComputerFilter: config.ComputerFilter,
		GroupFilter:    config.GroupFilter,
		Attributes: []string{"userPrincipalName", "dn", "cn", "company", "department", "title", "telephoneNumber",
			"otherTelephone", "mail", "pager", "msRTCSIP-PrimaryUserAddress", "url", "memberOf", "displayName",
			"description", "userPrincipalName", "employeeNumber", "pwdLastSet", "proxyAddresses", "userAccountControl", "distinguishedName", "lastLogonTimestamp"},
	}
	return client
}
