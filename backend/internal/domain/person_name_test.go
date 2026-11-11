package domain

import (
	"errors"
	"strings"
	"testing"
)

func TestNormalizePersonName(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{name: "colapsa espacios", input: "  roxana   karina ", want: "roxana karina"},
		{name: "acepta ñ", input: "Ñahui", want: "Ñahui"},
		{name: "acepta apóstrofo", input: "O'Connor", want: "O'Connor"},
		{name: "acepta guion", input: "Pérez-Luna", want: "Pérez-Luna"},
		{name: "rechaza dígitos", input: "R2D2", wantErr: true},
		{name: "rechaza emoji", input: "😀", wantErr: true},
		{name: "rechaza un carácter", input: "A", wantErr: true},
		{name: "rechaza 101 caracteres", input: strings.Repeat("a", 101), wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NormalizePersonName(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("NormalizePersonName(%q) error = %v, se esperaba error = %v", tt.input, err, tt.wantErr)
			}
			if tt.wantErr && !errors.Is(err, ErrInvalidName) {
				t.Errorf("se esperaba ErrInvalidName, obtuve %v", err)
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("NormalizePersonName(%q) = %q, se esperaba %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestComposeFullName(t *testing.T) {
	if got := ComposeFullName("Roxana Karina", "Delgado Huamán"); got != "Roxana Karina Delgado Huamán" {
		t.Errorf("ComposeFullName() = %q", got)
	}
}
