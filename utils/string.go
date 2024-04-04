package utils

func StringArrayContains(stringArr []string, s string) bool {
	for _, curStr := range stringArr {
		if curStr == s {
			return true
		}
	}
	return false
}
