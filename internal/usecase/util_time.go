package usecase

import "time"

func parseTicketDate(sDate string) string {
	currentDate := time.Now()
	dat, err := time.Parse("02.01.2006", sDate)
	if err != nil {
		return ""
	}
	if (currentDate.Day() == dat.Day() && currentDate.Month() == dat.Month() && currentDate.Year() == dat.Year()) || currentDate.After(dat) {
		return ""
	} else {
		return dat.Format(time.RFC3339)
	}
}
