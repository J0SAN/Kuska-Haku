package domain

import (
	"testing"
	"time"
)

func TestEstimateETA(t *testing.T) {
	cases := []struct {
		name     string
		meters   float64
		speedKmh float64
		want     time.Duration
	}{
		{"1000 m a 36 km/h (10 m/s)", 1000, 36, 100 * time.Second},
		{"distancia cero", 0, 30, 0},
		{"distancia negativa", -50, 30, 0},
		{"bus detenido usa velocidad promedio (20 km/h)", 1000, 0, 180 * time.Second},
		{"velocidad bajo el umbral usa promedio", 1000, 5, 180 * time.Second},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := EstimateETA(c.meters, c.speedKmh)
			diff := got - c.want
			if diff < 0 {
				diff = -diff
			}
			if diff > time.Second {
				t.Fatalf("EstimateETA(%v, %v) = %v, esperaba %v", c.meters, c.speedKmh, got, c.want)
			}
		})
	}
}

func TestClassifyStop(t *testing.T) {
	if passed, d := ClassifyStop(0.6, 0.3, 700); passed || d != 700 {
		t.Fatalf("paradero adelante: passed=%v d=%v", passed, d)
	}
	if passed, d := ClassifyStop(0.1, 0.3, 400); !passed || d != 0 {
		t.Fatalf("paradero atrás: passed=%v d=%v", passed, d)
	}
	if passed, d := ClassifyStop(0.299, 0.3, 5); passed || d != 0 {
		t.Fatalf("dentro de tolerancia: passed=%v d=%v", passed, d)
	}
}

func TestIsOffRoute(t *testing.T) {
	if IsOffRoute(40) {
		t.Error("40 m del trazado es ruido GPS normal, no desvío")
	}
	if !IsOffRoute(250) {
		t.Error("250 m del trazado debería ser desvío")
	}
}
