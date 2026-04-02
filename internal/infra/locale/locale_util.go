package locale

import (
	"golang.org/x/text/language"
)

var idToLocale map[string]language.Tag

func init() {
	idToLocale = make(map[string]language.Tag)
	// Simplified approximation of available locales, or standard language parsing
	// golang.org/x/text/language.Parse will check for validity anyway.
}

func IsAvailableLanguage(languageID string) bool {
	_, err := language.Parse(languageID)
	return err == nil
}
