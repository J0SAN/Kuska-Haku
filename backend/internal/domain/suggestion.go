package domain

import (
	"errors"
	"strings"
	"time"
	"unicode/utf8"
)

const (
	SuggestionStatusNew  = "new"
	SuggestionStatusRead = "read"

	minSuggestionRunes = 10
	maxSuggestionRunes = 1000
)

type Suggestion struct {
	ID        string    `json:"id"`
	UserID    string    `json:"userId"`
	UserName  string    `json:"userName"`
	Message   string    `json:"message"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"createdAt"`
}

// ValidateSuggestionMessage cuenta caracteres UTF-8 (runas) y no bytes.
func ValidateSuggestionMessage(message string) error {
	n := utf8.RuneCountInString(strings.TrimSpace(message))
	if n < minSuggestionRunes {
		return errors.New("la sugerencia debe tener al menos 10 caracteres")
	}
	if n > maxSuggestionRunes {
		return errors.New("la sugerencia no puede superar los 1000 caracteres")
	}
	return nil
}
