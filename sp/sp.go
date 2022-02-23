package sp

import (
	"github.com/koltyakov/gosip"
	"github.com/koltyakov/gosip/api"
	strategy "github.com/koltyakov/gosip/auth/ntlm"
	"log"
)

var sp *api.SP

//Init ...
func Init() {
	var err error
	sp, err = ConnectSP()
	if err != nil {
		log.Fatal(err)
	}
}

//ConnectSP ...
func ConnectSP() (*api.SP, error) {
	auth := &strategy.AuthCnfg{}
	auth.SiteURL = "http://nod2.brnv.rw"
	auth.Username = "portal-reader"
	auth.Password = "port@l123"
	auth.Domain = "brnv.rw"
	client := &gosip.SPClient{AuthCnfg: auth}
	sp := api.NewSP(client)
	return sp, nil
}

//GetSP ...
func GetSP() *api.SP {
	return sp
}
