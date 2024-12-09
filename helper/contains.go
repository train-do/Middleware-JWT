package helper

import "strconv"

func Contains(slice []string, str string) bool {
	for _, item := range slice {
		if item == str {
			return true
		}
	}
	return false
}

func StringToBool(str string) bool {
	convBool, _ := strconv.ParseBool(str)
	return convBool
}

func StringToInt(num string) int {
	convInt, _ := strconv.Atoi(num)
	return convInt
}
func IntToString(num int) string {
	convStr := strconv.Itoa(num)
	return convStr
}