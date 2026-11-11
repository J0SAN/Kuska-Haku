package ports

import (
	"context"

	"kuska-haku-backend/internal/domain"
)

// === FASE 0: Observabilidad y Salud de Infraestructura ===

type HealthPinger interface {
	Ping(ctx context.Context) error
	Name() string
}

// === FASE 1: Tiempo Real y Posiciones GPS (MQTT / Redis / WebSocket) ===

type PositionPublisher interface {
	Publish(ctx context.Context, pos domain.Position) error
}

type PositionStore interface {
	SaveLastPosition(ctx context.Context, pos domain.Position) error
	GetLastPosition(ctx context.Context, busID string) (*domain.Position, error)
	ListActiveBusIDs(ctx context.Context, routeID string) ([]string, error)
}

type PositionBroker interface {
	PublishToRoute(ctx context.Context, routeID string, pos domain.Position) error
	Subscribe(ctx context.Context, routeID string) (<-chan domain.Position, func(), error)
}

// === FASE 2: Autenticación, Criptografía y Usuarios ===

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	GetByID(ctx context.Context, id string) (*domain.User, error)
	ListByRole(ctx context.Context, role domain.Role) ([]domain.User, error)
}

type PasswordHasher interface {
	Hash(password string) (string, error)
	Compare(hash, password string) error
}

type TokenService interface {
	Generate(user domain.User) (string, error)
	Parse(token string) (*TokenClaims, error)
}

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
// y mide la desviación en metros del bus respecto al trazado de la ruta.
type RouteProgressRepository interface {
	Progress(ctx context.Context, routeID string, from domain.Coordinate) (*domain.BusProgress, error)
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

// === FASE 5: Ingesta Autenticada de GPS e Historial de Recorridos ===

type PositionHistoryRepository interface {
	SaveBatch(ctx context.Context, positions []domain.Position) (inserted int, err error)
}

type PositionIngestor interface {
	Ingest(ctx context.Context, batch []domain.Position) (domain.IngestResult, error)
}

// === FASE 7: QR por vehículo y calificaciones de pasajeros ===

type VehicleQRRepository interface {
	ListWithFleet(ctx context.Context) ([]domain.VehicleQRSummary, error)
	Create(ctx context.Context, fleetID, token, createdBy string) (*domain.VehicleQR, error)
	Rotate(ctx context.Context, fleetID, token, createdBy string) (*domain.VehicleQR, error) // una transacción
	Revoke(ctx context.Context, fleetID string) error
	ResolveActive(ctx context.Context, token string) (*domain.VehicleQRTarget, error)
}

type VehicleFeedbackRepository interface {
	Create(ctx context.Context, feedback *domain.VehicleFeedback) error
	List(ctx context.Context, filter domain.VehicleFeedbackFilter) ([]domain.VehicleFeedback, error)
	Summaries(ctx context.Context, fromDate string) ([]domain.VehicleRatingSummary, error)
	UpdateStatus(ctx context.Context, id, status string) error
}

// QRTokenGenerator produce tokens opacos de 22 caracteres base64url (128 bits).
type QRTokenGenerator interface {
	New() (string, error)
}
