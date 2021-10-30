package utils

func ContainsStr(slc *[]string, str *string) bool {
	for _, s := range *slc {
		if s == *str {
			return true
		}
	}
	return false
}
