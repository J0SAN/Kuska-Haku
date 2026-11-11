package auth

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"kuska-haku-backend/internal/domain"
)

// QRTokenGenerator produce tokens opacos de 128 bits con el alfabeto base64url.
// El token no lleva información del vehículo: solo lo identifica en la base.
type QRTokenGenerator struct{}

func NewQRTokenGenerator() *QRTokenGenerator {
	return &QRTokenGenerator{}
}

// New devuelve 16 bytes aleatorios en base64url sin relleno (22 caracteres).
func (g *QRTokenGenerator) New() (string, error) {
	raw := make([]byte, 16)
	if _, err := rand.Read(raw); err != nil {
		return "", fmt.Errorf("auth: error generando token QR: %w", err)
	}
	token := base64.RawURLEncoding.EncodeToString(raw)
	if len(token) != domain.QRTokenLength {
		return "", fmt.Errorf("auth: token QR con longitud inesperada (%d)", len(token))
	}
	return token, nil
}
