// Package domain contiene las entidades y reglas de negocio puras de Kuska Haku.
//
// Regla estricta de Arquitectura Hexagonal: este paquete NO importa nada de
// internal/adapters, ni de ningún framework HTTP, driver SQL o librería MQTT.
// Solo depende de la librería estándar de Go para permitir pruebas unitarias
// aisladas sin levantar infraestructura externa.
package domain
