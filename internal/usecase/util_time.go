package usecase

import "time"

func parseTicketDate(sDate string) string {
	location, err := time.LoadLocation("Local")
	if err != nil {
		return ""
	}
	_ = location
	currentDate := time.Now()
	dat, err := time.ParseInLocation("02.01.2006", sDate, location)
	if err != nil {
		return ""
	}
	if (currentDate.Day() == dat.Day() && currentDate.Month() == dat.Month() && currentDate.Year() == dat.Year()) || currentDate.After(dat) {
		return ""
	} else {
		return dat.Format(time.RFC3339)
	}
}
