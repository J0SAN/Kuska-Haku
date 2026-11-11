package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"

	"kuska-haku-backend/internal/adapters/auth"
	httpAdapter "kuska-haku-backend/internal/adapters/http"
	"kuska-haku-backend/internal/adapters/mqtt"
	"kuska-haku-backend/internal/adapters/postgres"
	"kuska-haku-backend/internal/adapters/redis"
	"kuska-haku-backend/internal/adapters/websocket"
	"kuska-haku-backend/internal/domain"
	"kuska-haku-backend/internal/migrations"
	"kuska-haku-backend/internal/usecases"
)

func main() {
	// === 1. INFRAESTRUCTURA POSTGRESQL + POSTGIS ===
	dbConnStr := cleanConnectionURL(os.Getenv("DATABASE_URL"))
	if dbConnStr == "" {
		dbConnStr = "postgres://postgres:postgres@localhost:5432/kuskahaku?sslmode=disable"
	}

	// Los errores de conexión pueden incluir la URL completa: se ocultan las
	// contraseñas antes de escribirlos en el log.
	db, err := sql.Open("postgres", dbConnStr)
	if err != nil {
		log.Fatalf("Error al conectar a la base de datos: %s", redactSecrets(err.Error()))
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("La base de datos no responde: %s", redactSecrets(err.Error()))
	}

	if err := migrations.Run(db); err != nil {
		log.Fatalf("Fallo crítico al inicializar la base de datos: %s", redactSecrets(err.Error()))
	}

	// --- Inicialización de Repositorios ---
	userRepo := postgres.NewUserRepository(db)
	routeRepo := postgres.NewRouteRepository(db)
	stopRepo := postgres.NewStopRepository(db)
	fleetRepo := postgres.NewFleetRepository(db)
	incomeRepo := postgres.NewIncomeRepository(db)
	suggestionRepo := postgres.NewSuggestionRepository(db)
	statsRepo := postgres.NewStatsRepository(db)
	historyRepo := postgres.NewPositionHistoryRepository(db)

	// === 2. INFRAESTRUCTURA TIEMPO REAL E INGESTA UNIFICADA (Redis + Caso de Uso + MQTT) ===
	// REDIS_URL (redis://…, la usa Render Key Value) tiene prioridad sobre REDIS_ADDR (host:puerto).
	var redisClient *redis.RedisPublisher
	if redisURL := cleanConnectionURL(os.Getenv("REDIS_URL")); redisURL != "" {
		redisClient, err = redis.NewRedisPublisherFromURL(redisURL)
		if err != nil {
			log.Fatalf("Configuración de Redis inválida: %s", redactSecrets(err.Error()))
		}
	} else {
		redisAddr := os.Getenv("REDIS_ADDR")
		if redisAddr == "" {
			redisAddr = "127.0.0.1:6379"
		}
		redisClient = redis.NewRedisPublisher(redisAddr)
	}
	ingestor := usecases.NewPositionIngestor(redisClient, redisClient, historyRepo)

	// MQTT es opcional: lo usan el simulador y los rastreadores GPS. La app del
	// conductor envía por HTTP, así que en despliegues sin broker se apaga con
	// MQTT_ENABLED=false.
	if os.Getenv("MQTT_ENABLED") == "false" {
		log.Println("MQTT desactivado (MQTT_ENABLED=false): solo se reciben posiciones por HTTP")
	} else {
		mqttBroker := os.Getenv("MQTT_BROKER_URL")
		if mqttBroker == "" {
			mqttBroker = "tcp://127.0.0.1:1883"
		}
		mqttAdapter := mqtt.NewMqttAdapter(mqttBroker, ingestor)
		mqttAdapter.Start()
	}

	// --- Inicialización de Servicios de Seguridad ---
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Println("ADVERTENCIA: JWT_SECRET no definido, se usa la clave de desarrollo")
		jwtSecret = "super-secret-key-para-desarrollo"
	}

	hasherSvc := auth.NewBcryptService()
	tokenSvc := auth.NewJWTService(jwtSecret)
	passengerRegistrar := usecases.NewPassengerRegistrar(userRepo, hasherSvc, tokenSvc)

	// --- Inicialización de Gateway WebSocket (con JWT) y Controladores HTTP ---
	wsGateway := websocket.NewGateway(redisClient, redisClient, tokenSvc)

	authController := httpAdapter.NewAuthController(userRepo, hasherSvc, tokenSvc, passengerRegistrar)
	// Detrás del proxy de Render todas las peticiones llegan desde su IP; con
	// CLIENT_IP_HEADER=CF-Connecting-IP cada cliente tiene su propia cuota.
	clientIPHeader := os.Getenv("CLIENT_IP_HEADER")
	registerLimiter := httpAdapter.NewRateLimiter(5, time.Hour).WithClientIPHeader(clientIPHeader)
	loginLimiter := httpAdapter.NewRateLimiter(10, 10*time.Minute).WithClientIPHeader(clientIPHeader)
	authMiddleware := httpAdapter.NewAuthMiddleware(tokenSvc)
	adminController := httpAdapter.NewAdminController(routeRepo, fleetRepo, stopRepo)
	passengerController := httpAdapter.NewPassengerController(routeRepo, fleetRepo, redisClient, routeRepo)
	driverController := httpAdapter.NewDriverController(fleetRepo, routeRepo, incomeRepo, ingestor)
	suggestionController := httpAdapter.NewSuggestionController(suggestionRepo)

	// === 3. ENRUTADOR (MUX) ===
	mux := http.NewServeMux()

	// Salud y Tiempo Real (WebSocket autenticado por ruta)
	mux.HandleFunc("/health", httpAdapter.HealthCheckHandler)
	mux.HandleFunc("/ws/routes/", wsGateway.ServeRoute)
	mux.HandleFunc("/ws", wsGateway.ServeRoute)

	// Autenticación (Fase 2)
	mux.HandleFunc("POST /api/v1/auth/register", registerLimiter.Wrap(authController.Register))
	mux.HandleFunc("POST /api/v1/auth/login", loginLimiter.Wrap(authController.Login))
	mux.HandleFunc("GET /api/v1/auth/me", authMiddleware.RequireAuth(authController.Me))

	// Panel Administrativo: Rutas, Paraderos, Flota y Usuarios (Fase 3, 4 & 5 - Solo Admin)
	mux.HandleFunc("POST /api/v1/admin/routes", authMiddleware.RequireRole(adminController.CreateRoute, domain.RoleAdmin))
	mux.HandleFunc("GET /api/v1/admin/routes", authMiddleware.RequireRole(adminController.GetRoutes, domain.RoleAdmin))
	mux.HandleFunc("/api/v1/admin/routes/", authMiddleware.RequireRole(adminController.RouteItem, domain.RoleAdmin))
	mux.HandleFunc("DELETE /api/v1/admin/stops/", authMiddleware.RequireRole(adminController.DeleteStop, domain.RoleAdmin))

	mux.HandleFunc("POST /api/v1/admin/fleet", authMiddleware.RequireRole(adminController.CreateFleet, domain.RoleAdmin))
	mux.HandleFunc("GET /api/v1/admin/fleet", authMiddleware.RequireRole(adminController.GetFleet, domain.RoleAdmin))
	mux.HandleFunc("/api/v1/admin/fleet/", authMiddleware.RequireRole(adminController.FleetItem, domain.RoleAdmin))
	mux.HandleFunc("GET /api/v1/admin/users", authMiddleware.RequireRole(httpAdapter.UsersHandler(userRepo), domain.RoleAdmin))

	// Pasajeros y consulta de ETA (Fase 4 & 5 - Cualquier usuario autenticado)
	mux.HandleFunc("GET /api/v1/routes", authMiddleware.RequireAuth(passengerController.ListRoutes))
	mux.HandleFunc("/api/v1/routes/", authMiddleware.RequireAuth(passengerController.RouteItem))

	// Conductor: Mi Bus, Ingresos e Ingesta de Posiciones GPS (Fase 4 & 5 - Solo Driver)
	mux.HandleFunc("GET /api/v1/driver/my-bus", authMiddleware.RequireRole(driverController.MyBus, domain.RoleDriver))
	mux.HandleFunc("/api/v1/driver/incomes", authMiddleware.RequireRole(driverController.Incomes, domain.RoleDriver))
	mux.HandleFunc("POST /api/v1/driver/positions", authMiddleware.RequireRole(driverController.Positions, domain.RoleDriver))

	// Buzón de Sugerencias (Fase 4)
	mux.HandleFunc("POST /api/v1/suggestions", authMiddleware.RequireRole(suggestionController.Create, domain.RolePassenger, domain.RoleDriver))
	mux.HandleFunc("GET /api/v1/admin/suggestions", authMiddleware.RequireRole(suggestionController.AdminCollection, domain.RoleAdmin))
	mux.HandleFunc("/api/v1/admin/suggestions/", authMiddleware.RequireRole(suggestionController.AdminItem, domain.RoleAdmin))

	// Dashboard Operativo (Fase 4 - Solo Admin)
	mux.HandleFunc("GET /api/v1/admin/dashboard", authMiddleware.RequireRole(httpAdapter.DashboardHandler(statsRepo), domain.RoleAdmin))

	// === 4. MIDDLEWARE CORS & SERVIDOR ===
	allowedOrigins := os.Getenv("CORS_ALLOWED_ORIGINS")
	if allowedOrigins == "" {
		allowedOrigins = httpAdapter.DefaultAllowedOrigins
	}
	handlerConCORS := httpAdapter.CORSMiddleware(httpAdapter.ParseAllowedOrigins(allowedOrigins), mux)

	// Render inyecta PORT; en local se mantiene :8080.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("kuska-haku backend escuchando en :%s", port)
	if err := http.ListenAndServe(":"+port, handlerConCORS); err != nil {
		log.Fatalf("Error en el servidor: %v", err)
	}
}
