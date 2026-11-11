package usecases

import (
	"context"
	"time"

	"kuska-haku-backend/internal/domain"
	"kuska-haku-backend/internal/ports"
)

// RegisterInput contiene los datos que envía el formulario de registro público.
type RegisterInput struct {
	DNI         string
	CheckDigit  string
	FirstNames  string
	LastNames   string
	Email       string
	Password    string
	AcceptTerms bool
}

// PassengerRegistrar crea cuentas de pasajero con DNI. El registro público
// nunca crea conductores ni administradores.
type PassengerRegistrar struct {
	users  ports.UserRepository
	hasher ports.PasswordHasher
	tokens ports.TokenService
	now    func() time.Time // inyectable para pruebas deterministas
}

func NewPassengerRegistrar(
	users ports.UserRepository,
	hasher ports.PasswordHasher,
	tokens ports.TokenService,
) *PassengerRegistrar {
	return &PassengerRegistrar{
		users:  users,
		hasher: hasher,
		tokens: tokens,
		now:    time.Now,
	}
}

// Register valida los datos, crea el usuario y devuelve el JWT de la nueva sesión.
func (p *PassengerRegistrar) Register(ctx context.Context, in RegisterInput) (*domain.User, string, error) {
	if !in.AcceptTerms {
		return nil, "", domain.ErrTermsNotAccepted
	}

	dni, err := domain.ParseDNI(in.DNI)
	if err != nil {
		return nil, "", err
	}
	if !dni.ValidCheckDigit(in.CheckDigit) {
		return nil, "", domain.ErrInvalidCheckDigit
	}

	firstNames, err := domain.NormalizePersonName(in.FirstNames)
	if err != nil {
		return nil, "", err
	}
	lastNames, err := domain.NormalizePersonName(in.LastNames)
	if err != nil {
		return nil, "", err
	}

	email := domain.NormalizeEmail(in.Email)
	if err := domain.ValidateEmail(email); err != nil {
		return nil, "", err
	}

	if len(in.Password) < 8 {
		return nil, "", domain.ErrWeakPassword
	}

	hashedPassword, err := p.hasher.Hash(in.Password)
	if err != nil {
		return nil, "", err
	}

	dniValue := string(dni)
	acceptedAt := p.now()
	user := &domain.User{
		Email:           email,
		PasswordHash:    hashedPassword,
		FullName:        domain.ComposeFullName(firstNames, lastNames),
		Role:            string(domain.RolePassenger),
		DNI:             &dniValue,
		TermsAcceptedAt: &acceptedAt,
	}

	if err := p.users.Create(ctx, user); err != nil {
		return nil, "", err
	}

	token, err := p.tokens.Generate(*user)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}
