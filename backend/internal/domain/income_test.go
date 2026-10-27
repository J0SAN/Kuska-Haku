package domain

import (
	"strings"
	"testing"
)

func TestIncomeRecord_Validate(t *testing.T) {
	ok := IncomeRecord{Passengers: 25, Amount: 62.5}
	if err := ok.Validate(); err != nil {
		t.Fatalf("registro válido rechazado: %v", err)
	}
	onlyAmount := IncomeRecord{Amount: 10}
	if err := onlyAmount.Validate(); err != nil {
		t.Fatalf("un registro solo con monto debería ser válido: %v", err)
	}

	cases := []struct {
		name string
		r    IncomeRecord
	}{
		{"pasajeros negativos", IncomeRecord{Passengers: -1, Amount: 5}},
		{"monto negativo", IncomeRecord{Passengers: 5, Amount: -1}},
		{"tres decimales", IncomeRecord{Passengers: 5, Amount: 1.234}},
		{"todo en cero", IncomeRecord{}},
		{"monto absurdo", IncomeRecord{Passengers: 1, Amount: 1e9}},
		{"pasajeros absurdos", IncomeRecord{Passengers: 100000, Amount: 5}},
		{"notas demasiado largas", IncomeRecord{Passengers: 1, Notes: strings.Repeat("a", 301)}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if err := c.r.Validate(); err == nil {
				t.Fatalf("esperaba error para: %s", c.name)
			}
		})
	}
}

func TestTodayLima_Format(t *testing.T) {
	d := TodayLima()
	if len(d) != 10 || d[4] != '-' || d[7] != '-' {
		t.Fatalf("formato inesperado: %q", d)
	}
}
