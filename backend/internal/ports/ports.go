package ports

import (
	"context"

	"kuska-haku-backend/internal/domain"
)

// === FASE 0: Observabilidad y Salud de Infraestructura ===

// HealthPinger lo implementa cualquier dependencia (Postgres, Redis, MQTT)
// capaz de reportar su estado de conexión al endpoint /health.
type HealthPinger interface {
	Ping(ctx context.Context) error
	Name() string
}

// === FASE 1: Tiempo Real y Posiciones GPS (MQTT / Redis / WebSocket) ===

// PositionPublisher publica posiciones en tiempo real hacia el broker (Pub/Sub).
type PositionPublisher interface {
	Publish(ctx context.Context, pos domain.Position) error
}

// PositionStore almacena y consulta la última posición GPS conocida de cada bus en Redis.
type PositionStore interface {
	SaveLastPosition(ctx context.Context, pos domain.Position) error
	GetLastPosition(ctx context.Context, busID string) (*domain.Position, error)
	ListActiveBusIDs(ctx context.Context, routeID string) ([]string, error)
}

// PositionBroker desacopla al consumidor MQTT del gateway WebSocket mediante canales por ruta.
type PositionBroker interface {
	PublishToRoute(ctx context.Context, routeID string, pos domain.Position) error
	Subscribe(ctx context.Context, routeID string) (<-chan domain.Position, func(), error)
}

// === FASE 2: Autenticación, Criptografía y Usuarios (Segregados por Responsabilidad) ===

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	GetByID(ctx context.Context, id string) (*domain.User, error)
}

// PasswordHasher desacopla el algoritmo de hashing (bcrypt hoy, Argon2id a futuro) del dominio.
type PasswordHasher interface {
	Hash(password string) (string, error)
	Compare(hash, password string) error
}

// TokenService emite y valida tokens criptográficos de sesión (JWT).
type TokenService interface {
	Generate(user domain.User) (string, error)
	Parse(token string) (*TokenClaims, error)
}

// TokenClaims contiene únicamente la identidad verificada extraída del JWT en cada request,
// evitando consultar la base de datos o exponer campos sensibles como PasswordHash.
type TokenClaims struct {
	UserID string
	Email  string
	Role   domain.Role
}

// === FASE 3: Flota, Paraderos y Trazado Espacial (PostGIS) ===

type RouteRepository interface {
	Create(ctx context.Context, route *domain.Route) error
	GetAll(ctx context.Context) ([]domain.Route, error)
	GetByID(ctx context.Context, id string) (*domain.Route, error)
	Update(ctx context.Context, route *domain.Route) error
	Delete(ctx context.Context, id string) error
	CountBuses(ctx context.Context, routeID string) (int, error)
}

type StopRepository interface {
	Create(ctx context.Context, stop *domain.Stop) error
	GetByRouteID(ctx context.Context, routeID string) ([]domain.Stop, error)
	Delete(ctx context.Context, id string) error
}

type FleetRepository interface {
	Create(ctx context.Context, fleet *domain.Fleet) error
	GetAll(ctx context.Context) ([]domain.Fleet, error)
	GetByID(ctx context.Context, id string) (*domain.Fleet, error)
	Update(ctx context.Context, fleet *domain.Fleet) error
	Delete(ctx context.Context, id string) error
	GetByRouteID(ctx context.Context, routeID string) ([]domain.Fleet, error)
	GetByDriverID(ctx context.Context, driverID string) (*domain.Fleet, error)
}

// === FASE 4: ETA, Ingresos del Conductor, Sugerencias y Dashboard ===

// RouteProgressRepository calcula distancias sobre el LineString en PostGIS
// (usando ST_LineLocatePoint y ST_LineSubstring).
type RouteProgressRepository interface {
	DistancesToStops(ctx context.Context, routeID string, from domain.Coordinate) ([]domain.StopDistance, error)
}

type IncomeRepository interface {
	Create(ctx context.Context, record *domain.IncomeRecord) error
	GetByDriverAndDate(ctx context.Context, driverID, date string) ([]domain.IncomeRecord, error)
}

type SuggestionRepository interface {
	Create(ctx context.Context, suggestion *domain.Suggestion) error
	GetAll(ctx context.Context, onlyNew bool) ([]domain.Suggestion, error)
	MarkRead(ctx context.Context, id string) error
}

type StatsRepository interface {
	GetDashboard(ctx context.Context, today string, days int) (*domain.Dashboard, error)
}
