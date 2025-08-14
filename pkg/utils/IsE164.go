package utils

func IsE164(s string) bool {
	if len(s) < 4 || len(s) > 17 {
		return false
	}
	if s[0] != '+' {
		return false
	}
	for i := 1; i < len(s); i++ {
		if s[i] < '0' || s[i] > '9' {
			return false
		}
	}
	return true
}
