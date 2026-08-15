package models_adds

import "regexp"

// HasURL проверяет, есть ли в тексте ссылка (URL).
func HasURL(text string) bool {
	urlPattern := `(https?://)?(www\.)?[a-zA-Z0-9-]+\.[a-zA-Z]{2,}(/[^ \n]*)?`
	re := regexp.MustCompile(urlPattern)

	return re.MatchString(text)
}
