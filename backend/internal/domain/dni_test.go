package domain

import (
	"errors"
	"testing"
)

func TestParseDNI(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    DNI
		wantErr bool
	}{
		{name: "recorta espacios", input: " 46027897 ", want: DNI("46027897")},
		{name: "siete dígitos", input: "4602789", wantErr: true},
		{name: "nueve dígitos", input: "460278971", wantErr: true},
		{name: "letra al final", input: "4602789A", wantErr: true},
		{name: "vacío", input: "", wantErr: true},
		{name: "ancho completo", input: "４６０２７８９７", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseDNI(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ParseDNI(%q) error = %v, se esperaba error = %v", tt.input, err, tt.wantErr)
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("ParseDNI(%q) = %q, se esperaba %q", tt.input, got, tt.want)
			}
			if tt.wantErr && !errors.Is(err, ErrInvalidDNI) {
				t.Errorf("se esperaba ErrInvalidDNI, obtuve %v", err)
			}
		})
	}
}

// La tabla usa DNIs ficticios calculados con el mismo algoritmo. No se
// commitean DNIs reales (spec 002, §8).
func TestValidCheckDigit(t *testing.T) {
	tests := []struct {
		name       string
		dni        DNI
		checkDigit string
		want       bool
	}{
		{name: "12345678 con número", dni: "12345678", checkDigit: "1", want: true},
		{name: "12345678 con letra", dni: "12345678", checkDigit: "E", want: true},
		{name: "12345678 con letra minúscula", dni: "12345678", checkDigit: "e", want: true},
		{name: "12345678 con número incorrecto", dni: "12345678", checkDigit: "2", want: false},
		{name: "87654321 con número", dni: "87654321", checkDigit: "0", want: true},
		{name: "87654321 con letra", dni: "87654321", checkDigit: "D", want: true},
		{name: "00000001 con número", dni: "00000001", checkDigit: "4", want: true},
		{name: "11111111 con número", dni: "11111111", checkDigit: "7", want: true},
		{name: "11111111 con letra minúscula", dni: "11111111", checkDigit: "a", want: true},
		{name: "99999999 con letra", dni: "99999999", checkDigit: "I", want: true},
		{name: "dígito vacío", dni: "12345678", checkDigit: "", want: false},
		{name: "dígito de dos caracteres", dni: "12345678", checkDigit: "11", want: false},
		{name: "DNI inválido", dni: "1234", checkDigit: "1", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.dni.ValidCheckDigit(tt.checkDigit); got != tt.want {
				t.Errorf("ValidCheckDigit(%q) = %v, se esperaba %v", tt.checkDigit, got, tt.want)
			}
		})
	}
}

func TestDNIMasked(t *testing.T) {
	if got := DNI("46027897").Masked(); got != "****7897" {
		t.Errorf("Masked() = %q, se esperaba ****7897", got)
	}
	if got := DNI("123").Masked(); got != "****" {
		t.Errorf("Masked() con DNI corto = %q, se esperaba ****", got)
	}
}
