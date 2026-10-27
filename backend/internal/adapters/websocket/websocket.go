package websocket

import (
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"

	"kuska-haku-backend/internal/ports"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type Gateway struct {
	store  ports.PositionStore
	broker ports.PositionBroker
	tokens ports.TokenService
}

func NewGateway(store ports.PositionStore, broker ports.PositionBroker, tokens ports.TokenService) *Gateway {
	return &Gateway{store: store, broker: broker, tokens: tokens}
}

// ServeRoute expone GET /ws/routes/{routeId}?token=<JWT>
func (g *Gateway) ServeRoute(w http.ResponseWriter, r *http.Request) {
	routeID := strings.TrimPrefix(r.URL.Path, "/ws/routes/")
	if routeID == "" || routeID == r.URL.Path {
		routeID = r.URL.Query().Get("routeId")
	}
	if routeID == "" {
		http.Error(w, "routeId requerido en la ruta, ej. /ws/routes/ruta-01", http.StatusBadRequest)
		return
	}

	if _, err := g.tokens.Parse(r.URL.Query().Get("token")); err != nil {
		http.Error(w, "token inválido o ausente", http.StatusUnauthorized)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("ws: error en upgrade para ruta %s: %v", routeID, err)
		return
	}
	defer conn.Close()

	ctx := r.Context()
	positions, unsubscribe, err := g.broker.Subscribe(ctx, routeID)
	if err != nil {
		log.Printf("ws: error suscribiendo a ruta %s: %v", routeID, err)
		return
	}
	defer unsubscribe()

	log.Printf("ws: cliente conectado a ruta %s", routeID)

	if busIDs, err := g.store.ListActiveBusIDs(ctx, routeID); err == nil {
		for _, busID := range busIDs {
			if pos, err := g.store.GetLastPosition(ctx, busID); err == nil && pos != nil {
				if err := conn.WriteJSON(pos); err != nil {
					return
				}
			}
		}
	}

	for pos := range positions {
		if err := conn.WriteJSON(pos); err != nil {
			log.Printf("ws: cliente de ruta %s desconectado: %v", routeID, err)
			return
		}
	}
}
