package domain

import "strings"

// DNI es el número de documento nacional de identidad peruano (8 dígitos).
type DNI string

// checkDigitWeights y las tablas de conversión implementan el algoritmo de
// dígito verificador de uso difundido para el DNI peruano. No está publicado
// oficialmente por RENIEC: solo comprueba que el número se copió bien, no que
// el documento pertenezca a quien lo escribe.
var (
	checkDigitWeights = [8]int{3, 2, 7, 6, 5, 4, 3, 2}
	checkDigitNumbers = [11]byte{'6', '7', '8', '9', '0', '1', '1', '2', '3', '4', '5'}
	checkDigitLetters = [11]byte{'K', 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J'}
)

// ParseDNI recorta espacios y exige exactamente 8 dígitos ASCII.
func ParseDNI(s string) (DNI, error) {
	s = strings.TrimSpace(s)
	if len(s) != 8 {
		return "", ErrInvalidDNI
	}
	for i := 0; i < len(s); i++ {
		if s[i] < '0' || s[i] > '9' {
			return "", ErrInvalidDNI
		}
	}
	return DNI(s), nil
}

// ValidCheckDigit comprueba el carácter impreso junto al número. Acepta el
// dígito numérico o la letra equivalente, en mayúscula o minúscula.
func (d DNI) ValidCheckDigit(c string) bool {
	if len(d) != 8 {
		return false
	}
	if len(c) != 1 {
		return false
	}

	sum := 0
	for i := 0; i < 8; i++ {
		digit := int(d[i] - '0')
		if digit < 0 || digit > 9 {
			return false
		}
		sum += digit * checkDigitWeights[i]
	}

	index := 11 - (sum % 11)
	if index == 11 {
		index = 0
	}

	c = strings.ToUpper(c)
	return c[0] == checkDigitNumbers[index] || c[0] == checkDigitLetters[index]
}

// Masked devuelve el DNI listo para logs: "****7897".
func (d DNI) Masked() string {
	if len(d) < 4 {
		return "****"
	}
	return "****" + string(d[len(d)-4:])
}
