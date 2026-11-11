package domain

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

const (
	personNameMinRunes = 2
	personNameMaxRunes = 100
)

// NormalizePersonName recorta, colapsa espacios internos y valida un nombre o
// apellido: entre 2 y 100 caracteres, solo letras Unicode (incluye ñ y tildes),
// espacio, apóstrofo y guion. Rechaza dígitos y emojis.
func NormalizePersonName(s string) (string, error) {
	normalized := strings.Join(strings.Fields(s), " ")

	runes := utf8.RuneCountInString(normalized)
	if runes < personNameMinRunes || runes > personNameMaxRunes {
		return "", ErrInvalidName
	}

	for _, r := range normalized {
		if unicode.IsLetter(r) || r == ' ' || r == '\'' || r == '\u2019' || r == '-' {
			continue
		}
		return "", ErrInvalidName
	}

	return normalized, nil
}

// ComposeFullName arma el nombre completo que se guarda en users.full_name.
func ComposeFullName(firstNames, lastNames string) string {
	return strings.TrimSpace(firstNames + " " + lastNames)
}
