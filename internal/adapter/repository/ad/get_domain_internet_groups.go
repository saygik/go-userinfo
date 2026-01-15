package ad

import (
	"strings"

	"github.com/saygik/go-userinfo/internal/config"
	"github.com/saygik/go-userinfo/internal/entity"
)

func (r *Repository) GetDomainInternetGroups(domain string) entity.ADInternetGroups {
	if val, ok := r.adconfigs[domain]; ok {
		if val.InternetGroups == (config.ADInternetGroups{}) {
			return entity.ADInternetGroups{WhiteList: []string{}, Full: []string{}, Tech: []string{}}
		}
		whiteList := []string{}
		full := []string{}
		tech := []string{}
		if val.InternetGroups.WhiteList != "" {
			whiteList = strings.Split(val.InternetGroups.WhiteList, ";")
		}
		if val.InternetGroups.Full != "" {
			full = strings.Split(val.InternetGroups.Full, ";")
		}
		if val.InternetGroups.Tech != "" {
			tech = strings.Split(val.InternetGroups.Tech, ";")
		}

		return entity.ADInternetGroups{WhiteList: whiteList, Full: full, Tech: tech}
	} else {
		return entity.ADInternetGroups{WhiteList: []string{}, Full: []string{}, Tech: []string{}}
	}
}
func (r *Repository) GetDomainInternetGroupsDN(domain string) entity.ADInternetGroupsDN {
	if val, ok := r.adconfigs[domain]; ok {
		return entity.ADInternetGroupsDN{WhiteList: val.InternetGroupsDN.WhiteList, Full: val.InternetGroupsDN.Full, Tech: val.InternetGroupsDN.Tech}
	} else {
		return entity.ADInternetGroupsDN{WhiteList: "", Full: "", Tech: ""}
	}
}
