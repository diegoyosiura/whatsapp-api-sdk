package utils

import "regexp"

func HasURL(text string) bool {
	return regexp.MustCompile(`https?://\S+`).MatchString(text)
}
