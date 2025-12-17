package models_adds

import "regexp"

// HasURL check if has url in text
func HasURL(text string) bool {
	urlPattern := `(https?://)?(www\.)?[a-zA-Z0-9-]+\.[a-zA-Z]{2,}(/[^ \n]*)?`
	re := regexp.MustCompile(urlPattern)

	return re.MatchString(text)
}
