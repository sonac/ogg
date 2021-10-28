package utils

import "os"

func Getenv(key string, defVal string) string {
	res := os.Getenv(key)
	if len(res) == 0 {
		return defVal
	}
	return res
}
