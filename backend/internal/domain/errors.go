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
)
