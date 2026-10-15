package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq" // Driver de Postgres

	"kuska-haku-backend/internal/adapters/auth"
	httpAdapter "kuska-haku-backend/internal/adapters/http"
	"kuska-haku-backend/internal/adapters/mqtt"
	"kuska-haku-backend/internal/adapters/postgres"
	"kuska-haku-backend/internal/adapters/redis"
	"kuska-haku-backend/internal/adapters/websocket"
	"kuska-haku-backend/internal/migrations"
)

type compositeAuthService struct {
	*auth.BcryptService
	*auth.JWTService
}

func main() {
	// === 1. INFRAESTRUCTURA FASE 1 (Sockets, Redis, MQTT) ===
	redisPublisher := redis.NewRedisPublisher("localhost:6379")
	mqttAdapter := mqtt.NewMqttAdapter("tcp://localhost:1883", redisPublisher)
	go mqttAdapter.Start()

	//Instanciar y arrancar el suscriptor de Redis en segundo plano
	redisSubscriber := redis.NewRedisSubscriber("localhost:6379")
	go redisSubscriber.Start()

	// === 2. INFRAESTRUCTURA FASE 2 & 3 (Postgres, Auth, JWT, PostGIS) ===
	dbConnStr := "postgres://postgres:postgres@localhost:5432/kuskahaku?sslmode=disable"
	db, err := sql.Open("postgres", dbConnStr)
	if err != nil {
		log.Fatalf("Error al conectar a la base de datos: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("La base de datos no responde: %v", err)
	}

	if err := migrations.Run(db); err != nil {
		log.Fatalf("Fallo crítico al inicializar la base de datos: %v", err)
	}

	// --- Inicialización de Repositorios ---
	userRepo := postgres.NewUserRepository(db)
	routeRepo := postgres.NewRouteRepository(db)
	fleetRepo := postgres.NewFleetRepository(db)

	// --- Inicialización de Servicios ---
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "super-secret-key-para-desarrollo"
	}

	authSvc := &compositeAuthService{
		BcryptService: auth.NewBcryptService(),
		JWTService:    auth.NewJWTService(jwtSecret),
	}

	// --- Inicialización de Controladores y Middlewares ---
	authController := httpAdapter.NewAuthController(userRepo, authSvc)
	authMiddleware := httpAdapter.NewAuthMiddleware(authSvc)
	adminController := httpAdapter.NewAdminController(routeRepo, fleetRepo)

	// === 3. ENRUTADOR (MUX) ===
	mux := http.NewServeMux()

	// Rutas Fase 1
	mux.HandleFunc("/health", httpAdapter.HealthCheckHandler)
	mux.HandleFunc("/ws", websocket.HandleWS)

	// Rutas Fase 2 (Públicas)
	mux.HandleFunc("POST /api/v1/auth/register", authController.Register)
	mux.HandleFunc("POST /api/v1/auth/login", authController.Login)

	// Rutas Fase 2 (Protegidas para validación de token)
	mux.HandleFunc("GET /api/v1/auth/me", authMiddleware.RequireAuth(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"message":"Acceso autorizado, token válido"}`))
	}))

	// Rutas Fase 3 (Protegidas - Exclusivas para Admin)
	// Envolvemos las funciones del controlador con el middleware RequireAuth
	mux.HandleFunc("POST /api/v1/admin/routes", authMiddleware.RequireAuth(adminController.CreateRoute))
	mux.HandleFunc("GET /api/v1/admin/routes", authMiddleware.RequireAuth(adminController.GetRoutes))
	mux.HandleFunc("POST /api/v1/admin/fleet", authMiddleware.RequireAuth(adminController.CreateFleet))
	mux.HandleFunc("GET /api/v1/admin/fleet", authMiddleware.RequireAuth(adminController.GetFleet))

	// === 4. MIDDLEWARE CORS ===
	// Envolvemos el multiplexor principal con el middleware para autorizar a React
	handlerConCORS := httpAdapter.CORSMiddleware(mux)

	// === 5. LEVANTAR EL SERVIDOR ===
	log.Println("kuska-haku backend escuchando en :8080")
	// Pasamos 'handlerConCORS' en lugar de 'mux'
	if err := http.ListenAndServe(":8080", handlerConCORS); err != nil {
		log.Fatalf("Error en el servidor: %v", err)
	}
}
