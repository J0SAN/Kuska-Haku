package domain

import "time"

const (
	// minSpeedKmh: por debajo de este umbral se considera que el bus está detenido
	// (semáforo, paradero, tráfico) y se usa avgSpeedKmh para evitar ETAs infinitos.
	minSpeedKmh = 8.0

	// avgSpeedKmh: velocidad urbana de referencia cuando el bus está detenido.
	avgSpeedKmh = 20.0

	// stopToleranceMeters: margen en metros para absorber el ruido del GPS cuando
	// el bus se encuentra detenido justo sobre el paradero.
	stopToleranceMeters = 15.0
)

// StopDistance representa la distancia medida a lo largo del trazado (LineString)
// desde la posición actual de un bus hasta un paradero específico.
type StopDistance struct {
	StopID         string  `json:"stopId"`
	Name           string  `json:"name"`
	StopOrder      int     `json:"stopOrder"`
	DistanceMeters float64 `json:"distanceMeters"`
	Passed         bool    `json:"passed"`
}

// ClassifyStop determina si un paradero ya fue superado por el bus a partir de
// sus fracciones relativas (0..1) sobre el LineString de la ruta.
func ClassifyStop(stopFraction, busFraction, distanceMeters float64) (passed bool, effectiveMeters float64) {
	if stopFraction < busFraction && distanceMeters > stopToleranceMeters {
		return true, 0
	}
	if stopFraction < busFraction {
		return false, 0 // dentro de la tolerancia: se considera en el paradero
	}
	return false, distanceMeters
}

// EstimateETA calcula la duración estimada para recorrer distanceMeters a speedKmh.
func EstimateETA(distanceMeters, speedKmh float64) time.Duration {
	if distanceMeters <= 0 {
		return 0
	}
	v := speedKmh
	if v < minSpeedKmh {
		v = avgSpeedKmh
	}
	metersPerSecond := v * 1000.0 / 3600.0
	seconds := distanceMeters / metersPerSecond
	return time.Duration(seconds * float64(time.Second))
}
