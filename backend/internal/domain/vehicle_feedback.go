package domain

import (
	"crypto/sha256"
	"encoding/hex"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"
)

const (
	QRTokenLength = 22

	FeedbackStatusNew    = "new"
	FeedbackStatusRead   = "read"
	FeedbackStatusHidden = "hidden"

	maxFeedbackCommentRunes = 500
	minRating               = 1
	maxRating               = 5
)

var (
	qrTokenPattern  = regexp.MustCompile(`^[A-Za-z0-9_-]{22}$`)
	deviceIDPattern = regexp.MustCompile(`^[a-f0-9]{32}$`)
)

// VehicleQR es un token opaco asociado a un vehículo. El token vive en su
// propia tabla y nunca se serializa como parte de domain.Fleet.
type VehicleQR struct {
	ID        string     `json:"id"`
	FleetID   string     `json:"fleetId"`
	Token     string     `json:"token"`
	CreatedBy *string    `json:"createdBy,omitempty"`
	CreatedAt time.Time  `json:"createdAt"`
	RevokedAt *time.Time `json:"revokedAt,omitempty"`
}

// VehicleQRTarget es lo que el token resuelve para el formulario público.
// Solo contiene datos del vehículo y la ruta: nunca del conductor.
type VehicleQRTarget struct {
	QRCodeID  string  `json:"-"`
	FleetID   string  `json:"-"`
	BusCode   string  `json:"busCode"`
	RouteName string  `json:"routeName"`
	DriverID  *string `json:"-"`
}

// VehicleQRSummary es una fila del listado administrativo: cada vehículo con
// su QR activo, si lo tiene.
type VehicleQRSummary struct {
	FleetID    string     `json:"fleetId"`
	BusCode    string     `json:"busCode"`
	RouteID    string     `json:"routeId"`
	RouteName  string     `json:"routeName"`
	DriverID   *string    `json:"driverId,omitempty"`
	DriverName string     `json:"driverName,omitempty"`
	QRCodeID   *string    `json:"qrCodeId,omitempty"`
	Token      *string    `json:"token,omitempty"`
	CreatedAt  *time.Time `json:"createdAt,omitempty"`
}

// VehicleFeedback es una calificación anónima de un vehículo concreto.
// DriverID es el conductor asignado en el momento de calificar; DriverName
// solo se rellena en el listado administrativo.
type VehicleFeedback struct {
	ID          string    `json:"id"`
	FleetID     *string   `json:"fleetId,omitempty"`
	QRCodeID    *string   `json:"qrCodeId,omitempty"`
	BusCode     string    `json:"busCode"`
	RouteName   string    `json:"routeName"`
	DriverID    *string   `json:"driverId,omitempty"`
	DriverName  string    `json:"driverName,omitempty"`
	Rating      int       `json:"rating"`
	Comment     string    `json:"comment"`
	DeviceHash  string    `json:"-"`
	ServiceDate string    `json:"serviceDate"`
	BusLive     bool      `json:"busLive"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"createdAt"`
}

// VehicleRatingSummary es el promedio por vehículo de un periodo.
type VehicleRatingSummary struct {
	FleetID   *string `json:"fleetId,omitempty"`
	BusCode   string  `json:"busCode"`
	RouteName string  `json:"routeName"`
	Average   float64 `json:"average"`
	Total     int     `json:"total"`
}

// VehicleFeedbackFilter acota el listado administrativo.
type VehicleFeedbackFilter struct {
	FleetID string
	Status  string
	Limit   int
}

// ValidateQRToken comprueba longitud y alfabeto base64url del token opaco.
func ValidateQRToken(token string) error {
	if !qrTokenPattern.MatchString(token) {
		return ErrQRNotAvailable
	}
	return nil
}

// ValidateRating exige una calificación entera entre 1 y 5.
func ValidateRating(rating int) error {
	if rating < minRating || rating > maxRating {
		return ErrInvalidRating
	}
	return nil
}

// NormalizeFeedbackComment recorta espacios y limita el comentario a 500 runas.
func NormalizeFeedbackComment(comment string) (string, error) {
	trimmed := strings.TrimSpace(comment)
	if utf8.RuneCountInString(trimmed) > maxFeedbackCommentRunes {
		return "", ErrInvalidFeedbackComment
	}
	return trimmed, nil
}

// ValidateDeviceID acepta 16 bytes en hexadecimal minúscula (32 caracteres).
func ValidateDeviceID(deviceID string) error {
	if !deviceIDPattern.MatchString(deviceID) {
		return ErrInvalidDeviceID
	}
	return nil
}

// HashDeviceID convierte el identificador del navegador en un hash SHA-256.
// La base nunca guarda el valor original.
func HashDeviceID(deviceID string) string {
	sum := sha256.Sum256([]byte(deviceID))
	return hex.EncodeToString(sum[:])
}

// ValidateFeedbackStatus comprueba el ciclo de vida administrativo.
func ValidateFeedbackStatus(status string) error {
	switch status {
	case FeedbackStatusNew, FeedbackStatusRead, FeedbackStatusHidden:
		return nil
	}
	return ErrInvalidFeedbackStatus
}
