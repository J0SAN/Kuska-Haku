package domain

import (
	"strings"
	"testing"
	"time"
)

func TestValidateRating(t *testing.T) {
	for _, valid := range []int{1, 5} {
		if err := ValidateRating(valid); err != nil {
			t.Fatalf("rating %d debe ser válido: %v", valid, err)
		}
	}
	for _, invalid := range []int{0, 6, -1} {
		if err := ValidateRating(invalid); err != ErrInvalidRating {
			t.Fatalf("rating %d debe rechazarse, obtuve %v", invalid, err)
		}
	}
}

func TestNormalizeFeedbackComment(t *testing.T) {
	got, err := NormalizeFeedbackComment("   ")
	if err != nil || got != "" {
		t.Fatalf("comentario vacío debe ser válido, obtuve %q / %v", got, err)
	}

	got, err = NormalizeFeedbackComment("  Buen servicio  ")
	if err != nil || got != "Buen servicio" {
		t.Fatalf("debe recortar espacios, obtuve %q / %v", got, err)
	}

	// 500 runas con tildes (cada tilde es 1 runa aunque ocupe 2 bytes).
	exact := strings.Repeat("á", 500)
	if got, err := NormalizeFeedbackComment(exact); err != nil || got != exact {
		t.Fatalf("500 runas deben aceptarse, obtuve %d runas / %v", len([]rune(got)), err)
	}

	if _, err := NormalizeFeedbackComment(strings.Repeat("á", 501)); err != ErrInvalidFeedbackComment {
		t.Fatalf("501 runas deben rechazarse, obtuve %v", err)
	}
}

func TestValidateQRToken(t *testing.T) {
	if err := ValidateQRToken("Abcdefghijklmnopqrstuv"); err != nil {
		t.Fatalf("token válido rechazado: %v", err)
	}
	for _, invalid := range []string{
		"",
		"corto",
		"Abcdefghijklmnopqrstu!",  // 22 con caracter fuera del alfabeto
		"Abcdefghijklmnopqrstuvw", // 23
		"Abcdefghijklmnopqrstu",   // 21
	} {
		if err := ValidateQRToken(invalid); err != ErrQRNotAvailable {
			t.Fatalf("token %q debe rechazarse, obtuve %v", invalid, err)
		}
	}
}

func TestValidateDeviceID(t *testing.T) {
	if err := ValidateDeviceID("0123456789abcdef0123456789abcdef"); err != nil {
		t.Fatalf("deviceID válido rechazado: %v", err)
	}
	for _, invalid := range []string{
		"",
		"0123456789ABCDEF0123456789abcdef", // mayúsculas
		"0123456789abcdef0123456789abcde",  // 31
		"0123456789abcdef0123456789abcdefg",
	} {
		if err := ValidateDeviceID(invalid); err != ErrInvalidDeviceID {
			t.Fatalf("deviceID %q debe rechazarse, obtuve %v", invalid, err)
		}
	}
}

func TestHashDeviceIDEsEstableYDe64Caracteres(t *testing.T) {
	a := HashDeviceID("0123456789abcdef0123456789abcdef")
	b := HashDeviceID("0123456789abcdef0123456789abcdef")
	if a != b || len(a) != 64 {
		t.Fatalf("hash inesperado: %q vs %q (len %d)", a, b, len(a))
	}
	if a == "0123456789abcdef0123456789abcdef" {
		t.Fatal("el hash no debe ser igual al dispositivo original")
	}
}

func TestValidateFeedbackStatus(t *testing.T) {
	for _, valid := range []string{FeedbackStatusNew, FeedbackStatusRead, FeedbackStatusHidden} {
		if err := ValidateFeedbackStatus(valid); err != nil {
			t.Fatalf("estado %q debe ser válido: %v", valid, err)
		}
	}
	if err := ValidateFeedbackStatus("deleted"); err != ErrInvalidFeedbackStatus {
		t.Fatalf("estado desconocido debe rechazarse, obtuve %v", err)
	}
}

func TestDateLima(t *testing.T) {
	// 03:00 UTC del 2 de enero son las 22:00 del 1 de enero en Lima (UTC-5).
	instant := time.Date(2026, time.January, 2, 3, 0, 0, 0, time.UTC)
	if got := DateLima(instant); got != "2026-01-01" {
		t.Fatalf("esperaba 2026-01-01, obtuve %s", got)
	}

	// 12:00 UTC sigue siendo el mismo día en Lima.
	instant = time.Date(2026, time.November, 11, 12, 0, 0, 0, time.UTC)
	if got := DateLima(instant); got != "2026-11-11" {
		t.Fatalf("esperaba 2026-11-11, obtuve %s", got)
	}
}
