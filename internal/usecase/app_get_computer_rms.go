package usecase

import (
	"fmt"

	"github.com/saygik/go-userinfo/internal/entity"
)

func (u *UseCase) GetComputerRMS(domain string) ([]entity.ComputerRMS, error) {
	computerRMS := u.ad.GetDomainGLPIOrganization(domain)
	if computerRMS == "" {
		return nil, fmt.Errorf("Нет организации GLPI для домена")
	}
	rms, err := u.repo.GetComputerRMS(domain)
	if err != nil {
		return nil, fmt.Errorf("Отсутствуют данные логон скриптоп по RMS")
	}

	versions, _ := u.glpi.GetAllRMSSoftwares(computerRMS)
	versionMap := make(map[string]entity.ComputerRMS, len(versions))
	for _, v := range versions {
		versionMap[v.Computer] = v
	}
	for i := range rms {
		if v, ok := versionMap[rms[i].Computer]; ok {
			rms[i].Version = v.Version
			rms[i].SoftDate = v.SoftDate
		}
	}
	return rms, nil
}
