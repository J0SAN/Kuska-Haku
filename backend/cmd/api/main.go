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

	// === 2. INFRAESTRUCTURA FASE 2 (Postgres, Auth, JWT) ===
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

	userRepo := postgres.NewUserRepository(db)

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "super-secret-key-para-desarrollo"
	}

	authSvc := &compositeAuthService{
		BcryptService: auth.NewBcryptService(),
		JWTService:    auth.NewJWTService(jwtSecret),
	}

	authController := httpAdapter.NewAuthController(userRepo, authSvc)
	authMiddleware := httpAdapter.NewAuthMiddleware(authSvc)

	// === 3. ENRUTADOR (MUX) ===
	mux := http.NewServeMux()

	// Rutas Fase 1
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	})
	mux.HandleFunc("/ws", websocket.HandleWS)

	// Rutas Fase 2 (Públicas)
	mux.HandleFunc("/api/v1/auth/register", authController.Register)
	mux.HandleFunc("/api/v1/auth/login", authController.Login)

	// Rutas Fase 2 (Protegidas)
	mux.HandleFunc("/api/v1/auth/me", authMiddleware.RequireAuth(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"message":"Acceso autorizado, token válido"}`))
	}))

	// === 4. LEVANTAR EL SERVIDOR ===
	log.Println("kuska-haku backend escuchando en :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Error en el servidor: %v", err)
	}
}
