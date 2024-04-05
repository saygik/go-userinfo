package usecase

import "github.com/saygik/go-userinfo/internal/entity"

func (u *UseCase) GetStatTop10Performers(startdate string, enddate string) ([]entity.GLPIStatsTop10, error) {
	return u.glpi.GetStatTop10Performers(startdate, enddate)
}

func (u *UseCase) GetStatTop10Iniciators(startdate string, enddate string) ([]entity.GLPIStatsTop10, error) {
	return u.glpi.GetStatTop10Iniciators(startdate, enddate)
}
func (u *UseCase) GetStatTop10Groups(startdate string, enddate string) ([]entity.GLPIStatsTop10, error) {
	return u.glpi.GetStatTop10Groups(startdate, enddate)
}

func (u *UseCase) GetStatTickets() ([]entity.GLPITicketsStats, error) {
	return u.glpi.GetStatTickets()
}

func (u *UseCase) GetStatTicketsDays(startdate string, enddate string) ([]entity.GLPITicketsStats, error) {
	return u.glpi.GetStatTicketsDays(startdate, enddate)
}
func (u *UseCase) GetStatRegions(startdate string, enddate string) (tickets []entity.GLPIRegionsStats, err error) {
	return u.glpi.GetStatRegions(startdate, enddate)
}
func (u *UseCase) GetStatPeriodRequestTypes(startdate string, enddate string) (tickets []entity.GLPIStatsTop10, err error) {
	return u.glpi.GetStatPeriodRequestTypes(startdate, enddate)
}
func (u *UseCase) GetStatPeriodRegionDayCounts(startdate string, enddate string, maxday int) (tickets []entity.RegionsDayStats, err error) {
	return u.glpi.GetStatPeriodRegionDayCounts(startdate, enddate, maxday)
}
func (u *UseCase) GetStatPeriodOrgTreemap(startdate string, enddate string) (tickets []entity.TreemapData, err error) {
	return u.glpi.GetStatPeriodOrgTreemap(startdate, enddate)
}
func (u *UseCase) GetStatFailures() ([]entity.GLPITicketsStats, error) {
	return u.glpi.GetStatFailures()
}
func (u *UseCase) GetStatPeriodTicketsCounts(startdate string, enddate string) ([]entity.GLPIStatsCounts, error) {
	return u.glpi.GetStatPeriodTicketsCounts(startdate, enddate)
}
