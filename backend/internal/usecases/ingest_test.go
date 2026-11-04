package usecases

import (
	"context"
	"errors"
	"testing"
	"time"

	"kuska-haku-backend/internal/domain"
	"kuska-haku-backend/internal/ports"
)

var fixedNow = time.Date(2026, 10, 1, 12, 0, 0, 0, time.UTC)

type fakeStore struct {
	ports.PositionStore
	last  map[string]domain.Position
	saved []domain.Position
}

func (f *fakeStore) GetLastPosition(_ context.Context, busID string) (*domain.Position, error) {
	if p, ok := f.last[busID]; ok {
		return &p, nil
	}
	return nil, nil
}

func (f *fakeStore) SaveLastPosition(_ context.Context, p domain.Position) error {
	f.saved = append(f.saved, p)
	return nil
}

type fakeBroker struct {
	ports.PositionBroker
	published []domain.Position
}

func (f *fakeBroker) PublishToRoute(_ context.Context, _ string, p domain.Position) error {
	f.published = append(f.published, p)
	return nil
}

type fakeHistory struct {
	got      []domain.Position
	existing int // cuántas filas ya existían (duplicados)
	err      error
}

func (f *fakeHistory) SaveBatch(_ context.Context, ps []domain.Position) (int, error) {
	if f.err != nil {
		return 0, f.err
	}
	f.got = ps
	return len(ps) - f.existing, nil
}

func makePos(secondsAgo int) domain.Position {
	return domain.Position{
		BusID:     "bus-01",
		RouteID:   "ruta-01",
		Latitude:  -7.16,
		Longitude: -78.51,
		SpeedKmh:  20,
		Heading:   90,
		Timestamp: fixedNow.Add(-time.Duration(secondsAgo) * time.Second),
	}
}

func newTestIngestor(store *fakeStore, br *fakeBroker, h *fakeHistory) *PositionIngestor {
	i := NewPositionIngestor(store, br, h)
	i.now = func() time.Time { return fixedNow }
	return i
}

func TestIngest_FreshPointSavedAndPublished(t *testing.T) {
	st, br, h := &fakeStore{}, &fakeBroker{}, &fakeHistory{}
	res, err := newTestIngestor(st, br, h).Ingest(context.Background(), []domain.Position{makePos(2)})
	if err != nil || res.Saved != 1 || !res.Live || len(br.published) != 1 {
		t.Fatalf("res=%+v err=%v published=%d", res, err, len(br.published))
	}
}

func TestIngest_DelayedBatchSavedButNotPublishedLive(t *testing.T) {
	st, br, h := &fakeStore{}, &fakeBroker{}, &fakeHistory{}
	res, _ := newTestIngestor(st, br, h).Ingest(context.Background(), []domain.Position{makePos(600), makePos(590), makePos(580)})
	if res.Saved != 3 {
		t.Errorf("debía guardar los 3 en el historial, guardó %d", res.Saved)
	}
	if res.Live || len(br.published) != 0 || len(st.saved) != 0 {
		t.Error("un lote de hace 10 minutos NO debe mostrarse como posición en vivo")
	}
}

func TestIngest_MixedBatchPublishesOnlyNewest(t *testing.T) {
	st, br, h := &fakeStore{}, &fakeBroker{}, &fakeHistory{}
	// Desordenado a propósito: el caso de uso debe ordenar cronológicamente.
	res, _ := newTestIngestor(st, br, h).Ingest(context.Background(), []domain.Position{makePos(300), makePos(3), makePos(200)})
	if len(br.published) != 1 || !br.published[0].Timestamp.Equal(fixedNow.Add(-3*time.Second)) {
		t.Fatalf("debía difundirse solo el punto de hace 3 s, published=%+v", br.published)
	}
	if res.Saved != 3 {
		t.Errorf("saved=%d, esperaba 3", res.Saved)
	}
	if !h.got[0].Timestamp.Before(h.got[2].Timestamp) {
		t.Error("el historial debe recibirse ordenado de más antiguo a más reciente")
	}
}

func TestIngest_DoesNotMoveBusBackwardsIfNewerExists(t *testing.T) {
	known := makePos(1) // más reciente que el que llega
	st := &fakeStore{last: map[string]domain.Position{"bus-01": known}}
	br := &fakeBroker{}
	res, _ := newTestIngestor(st, br, &fakeHistory{}).Ingest(context.Background(), []domain.Position{makePos(20)})
	if res.Live || len(br.published) != 0 {
		t.Error("no debe sobrescribir una posición en vivo más reciente con una más antigua")
	}
}

func TestIngest_DiscardsInvalidPointsIndividually(t *testing.T) {
	badSpeed := makePos(5)
	badSpeed.SpeedKmh = 900
	futurePoint := makePos(5)
	futurePoint.Timestamp = fixedNow.Add(time.Hour)

	st, br, h := &fakeStore{}, &fakeBroker{}, &fakeHistory{}
	res, _ := newTestIngestor(st, br, h).Ingest(context.Background(), []domain.Position{makePos(10), badSpeed, futurePoint, makePos(4)})
	if res.Received != 4 || res.Saved != 2 || res.Discarded != 2 {
		t.Fatalf("res=%+v", res)
	}
	if res.FirstDiscardReason == "" {
		t.Error("debe informar el motivo del primer descarte")
	}
}

func TestIngest_ReportsDuplicates(t *testing.T) {
	h := &fakeHistory{existing: 2}
	res, _ := newTestIngestor(&fakeStore{}, &fakeBroker{}, h).Ingest(context.Background(), []domain.Position{makePos(9), makePos(6), makePos(3)})
	if res.Saved != 1 || res.Duplicates != 2 {
		t.Fatalf("res=%+v", res)
	}
}

func TestIngest_HistoryErrorPreventsLivePublish(t *testing.T) {
	br := &fakeBroker{}
	h := &fakeHistory{err: errors.New("db caída")}
	_, err := newTestIngestor(&fakeStore{}, br, h).Ingest(context.Background(), []domain.Position{makePos(2)})
	if err == nil {
		t.Fatal("esperaba error cuando falla el historial")
	}
	if len(br.published) != 0 {
		t.Error("no debe difundir si no pudo guardar (el cliente reintentará el lote)")
	}
}
