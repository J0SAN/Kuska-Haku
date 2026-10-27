package domain

import (
	"strings"
	"testing"
)

func TestValidateSuggestionMessage(t *testing.T) {
	if err := ValidateSuggestionMessage("El bus pasa muy rápido por el mercado"); err != nil {
		t.Fatalf("mensaje válido rechazado: %v", err)
	}
	if err := ValidateSuggestionMessage("ñandúñandú"); err != nil {
		t.Fatalf("10 caracteres multibyte deberían ser válidos: %v", err)
	}
	if err := ValidateSuggestionMessage("hola"); err == nil {
		t.Fatal("un mensaje de 4 caracteres debería ser inválido")
	}
	if err := ValidateSuggestionMessage("          corto          "); err == nil {
		t.Fatal("los espacios laterales no deberían contar para el mínimo")
	}
	if err := ValidateSuggestionMessage(strings.Repeat("a", 1001)); err == nil {
		t.Fatal("más de 1000 caracteres debería ser inválido")
	}
}
