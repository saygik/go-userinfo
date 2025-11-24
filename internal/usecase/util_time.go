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

func parseDate(sDate string) string {
	parsedTime, err := time.Parse(time.RFC3339, sDate)
	if err != nil {
		return ""
	}
	return parsedTime.Format("02.01.2006")
}
func parseGlpiDate(sDate string) string {
	t, err := time.Parse("2006-01-02 15:04:05", sDate)
	if err != nil {
		return sDate
	}
	return t.Format("02.01.2006 15:04:05")
}
