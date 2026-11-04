// Package usecases contiene la lógica de aplicación: orquesta el dominio y
// los puertos, sin depender de detalles de infraestructura (Postgres, Redis o MQTT).
package usecases

import (
	"context"
	"fmt"
	"sort"
	"time"

	"kuska-haku-backend/internal/domain"
	"kuska-haku-backend/internal/ports"
)

// PositionIngestor implementa ports.PositionIngestor.
type PositionIngestor struct {
	store   ports.PositionStore
	broker  ports.PositionBroker
	history ports.PositionHistoryRepository
	now     func() time.Time // inyectable para pruebas deterministas
}

func NewPositionIngestor(
	store ports.PositionStore,
	broker ports.PositionBroker,
	history ports.PositionHistoryRepository,
) *PositionIngestor {
	return &PositionIngestor{
		store:   store,
		broker:  broker,
		history: history,
		now:     time.Now,
	}
}

// Ingest procesa un lote de posiciones GPS:
//  1. Descarta individualmente las inválidas (un punto defectuoso no aborta el lote).
//  2. Guarda todas las posiciones válidas en el historial de forma idempotente.
//  3. Por cada bus, difunde en vivo únicamente su punto más reciente, siempre que sea
//     fresco (< 60 s) y posterior a la última posición conocida en Redis.
func (i *PositionIngestor) Ingest(ctx context.Context, batch []domain.Position) (domain.IngestResult, error) {
	now := i.now()
	res := domain.IngestResult{Received: len(batch)}

	valid := make([]domain.Position, 0, len(batch))
	for _, p := range batch {
		if err := p.ValidateForIngest(now); err != nil {
			res.Discarded++
			if res.FirstDiscardReason == "" {
				res.FirstDiscardReason = err.Error()
			}
			continue
		}
		valid = append(valid, p)
	}
	if len(valid) == 0 {
		return res, nil
	}

	sort.Slice(valid, func(a, b int) bool {
		return valid[a].Timestamp.Before(valid[b].Timestamp)
	})

	inserted, err := i.history.SaveBatch(ctx, valid)
	if err != nil {
		return res, fmt.Errorf("ingest: error guardando el historial: %w", err)
	}
	res.Saved = inserted
	res.Duplicates = len(valid) - inserted

	// Selecciona el punto más reciente de cada bus (el slice ya está ordenado cronológicamente).
	newest := map[string]domain.Position{}
	for _, p := range valid {
		newest[p.BusID] = p
	}

	for busID, p := range newest {
		if !p.IsFresh(now) {
			continue // atrasado: se conserva en historial pero no se emite en vivo
		}
		last, err := i.store.GetLastPosition(ctx, busID)
		if err != nil {
			return res, fmt.Errorf("ingest: error leyendo la última posición: %w", err)
		}
		if last != nil && !p.Timestamp.After(last.Timestamp) {
			continue // ya existe una posición igual o más reciente en vivo
		}
		if err := i.store.SaveLastPosition(ctx, p); err != nil {
			return res, fmt.Errorf("ingest: error guardando la última posición: %w", err)
		}
		if err := i.broker.PublishToRoute(ctx, p.RouteID, p); err != nil {
			return res, fmt.Errorf("ingest: error difundiendo la posición: %w", err)
		}
		res.Live = true
	}
	return res, nil
}
