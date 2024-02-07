package usecase

func IsStringInArray(str string, arr interface{}) bool {
	if arr == nil {
		return false
	}
	for _, b := range arr.([]string) {
		if b == str {
			return true
		}
	}
	return false
}
