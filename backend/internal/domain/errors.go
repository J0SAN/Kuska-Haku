package domain

import "errors"

// Errores centinela del dominio de usuarios. Los adaptadores HTTP los mapean a
// códigos de estado concretos con errors.Is.
var (
	ErrInvalidDNI             = errors.New("el DNI debe tener exactamente 8 dígitos")
	ErrInvalidCheckDigit      = errors.New("el dígito verificador no coincide con el DNI")
	ErrInvalidName            = errors.New("el nombre solo admite letras, espacios, apóstrofos y guiones")
	ErrInvalidEmail           = errors.New("formato de correo inválido")
	ErrWeakPassword           = errors.New("la contraseña debe tener al menos 8 caracteres")
	ErrTermsNotAccepted       = errors.New("debes aceptar la política de privacidad")
	ErrDNIAlreadyRegistered   = errors.New("ya existe una cuenta con ese DNI")
	ErrEmailAlreadyRegistered = errors.New("ya existe una cuenta con ese correo")
	ErrRoleNotAllowed         = errors.New("el registro público solo crea pasajeros")

	// Fase 7: QR por vehículo y calificaciones de pasajeros.
	ErrQRNotAvailable           = errors.New("este código QR ya no está disponible")
	ErrQRAlreadyExists          = errors.New("el vehículo ya tiene un QR activo")
	ErrInvalidRating            = errors.New("la calificación debe estar entre 1 y 5 estrellas")
	ErrInvalidFeedbackComment   = errors.New("el comentario no puede superar los 500 caracteres")
	ErrInvalidDeviceID          = errors.New("el identificador del dispositivo es inválido")
	ErrFeedbackAlreadySubmitted = errors.New("ya calificaste este vehículo hoy")
	ErrInvalidFeedbackStatus    = errors.New("estado de calificación inválido")
)
