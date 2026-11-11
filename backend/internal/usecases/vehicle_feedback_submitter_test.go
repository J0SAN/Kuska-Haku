package usecases

import (
	"context"
	"errors"
	"testing"
	"time"

	"kuska-haku-backend/internal/domain"
	"kuska-haku-backend/internal/ports"
)

// --- Fakes que incrustan los puertos ---

type fakeQRRepo struct {
	ports.VehicleQRRepository
	target      *domain.VehicleQRTarget
	resolveErr  error
	createCalls int
	lastToken   string
	lastFleetID string
	lastUser    string
	rotateCalls int
	revokeCalls int
	createErr   error
}

func (f *fakeQRRepo) Create(_ context.Context, fleetID, token, createdBy string) (*domain.VehicleQR, error) {
	f.createCalls++
	f.lastFleetID, f.lastToken, f.lastUser = fleetID, token, createdBy
	if f.createErr != nil {
		return nil, f.createErr
	}
	return &domain.VehicleQR{ID: "qr-1", FleetID: fleetID, Token: token}, nil
}

func (f *fakeQRRepo) Rotate(_ context.Context, fleetID, token, createdBy string) (*domain.VehicleQR, error) {
	f.rotateCalls++
	f.lastFleetID, f.lastToken, f.lastUser = fleetID, token, createdBy
	return &domain.VehicleQR{ID: "qr-2", FleetID: fleetID, Token: token}, nil
}

func (f *fakeQRRepo) Revoke(_ context.Context, fleetID string) error {
	f.revokeCalls++
	f.lastFleetID = fleetID
	return nil
}

func (f *fakeQRRepo) ResolveActive(_ context.Context, _ string) (*domain.VehicleQRTarget, error) {
	return f.target, f.resolveErr
}

type fakeTokenGenerator struct {
	token string
	err   error
	calls int
}

func (f *fakeTokenGenerator) New() (string, error) {
	f.calls++
	return f.token, f.err
}

type fakeFeedbackRepo struct {
	ports.VehicleFeedbackRepository
	createCalls int
	created     *domain.VehicleFeedback
	createErr   error
}

func (f *fakeFeedbackRepo) Create(_ context.Context, feedback *domain.VehicleFeedback) error {
	f.createCalls++
	f.created = feedback
	if f.createErr != nil {
		return f.createErr
	}
	feedback.ID = "fb-1"
	feedback.Status = domain.FeedbackStatusNew
	return nil
}

type fakePositionStore struct {
	ports.PositionStore
	position *domain.Position
	err      error
}

func (f *fakePositionStore) GetLastPosition(_ context.Context, _ string) (*domain.Position, error) {
	return f.position, f.err
}

// --- Pruebas: manager de QR ---

func TestVehicleQRManager_GenerarCuandoYaExisteDevuelveError(t *testing.T) {
	repo := &fakeQRRepo{createErr: domain.ErrQRAlreadyExists}
	manager := NewVehicleQRManager(repo, &fakeTokenGenerator{token: "Abcdefghijklmnopqrstuv"})

	if _, err := manager.Generate(context.Background(), "fleet-1", "admin-1"); !errors.Is(err, domain.ErrQRAlreadyExists) {
		t.Fatalf("esperaba ErrQRAlreadyExists, obtuve %v", err)
	}
	if repo.lastToken != "Abcdefghijklmnopqrstuv" {
		t.Fatalf("el token generado no llegó al repositorio: %q", repo.lastToken)
	}
	if repo.lastUser != "admin-1" {
		t.Fatalf("esperaba createdBy admin-1, obtuve %q", repo.lastUser)
	}
}

func TestVehicleQRManager_RotarYRevocarLlamanAlPuerto(t *testing.T) {
	repo := &fakeQRRepo{}
	manager := NewVehicleQRManager(repo, &fakeTokenGenerator{token: "Abcdefghijklmnopqrstuv"})

	qr, err := manager.Rotate(context.Background(), "fleet-9", "admin-2")
	if err != nil {
		t.Fatalf("rotar falló: %v", err)
	}
	if repo.rotateCalls != 1 || qr.FleetID != "fleet-9" {
		t.Fatalf("rotar no llamó al puerto correctamente: %+v", repo)
	}

	if err := manager.Revoke(context.Background(), "fleet-9"); err != nil {
		t.Fatalf("revocar falló: %v", err)
	}
	if repo.revokeCalls != 1 || repo.lastFleetID != "fleet-9" {
		t.Fatalf("revocar no llamó al puerto correctamente: %+v", repo)
	}
}

// --- Pruebas: envío de calificaciones ---

func fbTarget() *domain.VehicleQRTarget {
	driver := "driver-1"
	return &domain.VehicleQRTarget{
		QRCodeID:  "qr-1",
		FleetID:   "fleet-1",
		BusCode:   "bus-01",
		RouteName: "R-01 Centro",
		DriverID:  &driver,
	}
}

func fbFixedNow() time.Time {
	// 03:00 UTC del 2 de enero = 22:00 del 1 de enero en Lima.
	return time.Date(2026, time.January, 2, 3, 0, 0, 0, time.UTC)
}

func newSubmitter(qr ports.VehicleQRRepository, feedback ports.VehicleFeedbackRepository, positions ports.PositionStore) *VehicleFeedbackSubmitter {
	s := NewVehicleFeedbackSubmitter(feedback, qr, positions)
	s.now = fbFixedNow
	return s
}

func fbInput() VehicleFeedbackInput {
	return VehicleFeedbackInput{
		Rating:   4,
		Comment:  "  Buen servicio  ",
		DeviceID: "0123456789abcdef0123456789abcdef",
	}
}

func TestVehicleFeedbackSubmitter_GuardaSnapshotsYHash(t *testing.T) {
	qrRepo := &fakeQRRepo{target: fbTarget()}
	feedbackRepo := &fakeFeedbackRepo{}
	submitter := newSubmitter(qrRepo, feedbackRepo, &fakePositionStore{})

	feedback, err := submitter.Submit(context.Background(), "Abcdefghijklmnopqrstuv", fbInput())
	if err != nil {
		t.Fatalf("envío falló: %v", err)
	}
	if feedbackRepo.createCalls != 1 {
		t.Fatalf("esperaba una escritura, obtuve %d", feedbackRepo.createCalls)
	}
	if feedback.BusCode != "bus-01" || feedback.RouteName != "R-01 Centro" {
		t.Fatalf("snapshots incorrectos: %+v", feedback)
	}
	if feedback.DriverID == nil || *feedback.DriverID != "driver-1" {
		t.Fatalf("debe conservar el conductor asignado: %+v", feedback.DriverID)
	}
	if feedback.Comment != "Buen servicio" {
		t.Fatalf("el comentario debe recortarse, obtuve %q", feedback.Comment)
	}
	if feedback.DeviceHash != domain.HashDeviceID(fbInput().DeviceID) {
		t.Fatalf("deviceHash incorrecto: %q", feedback.DeviceHash)
	}
	if feedback.DeviceHash == fbInput().DeviceID {
		t.Fatal("nunca debe guardarse el identificador original")
	}
	if feedback.ServiceDate != "2026-01-01" {
		t.Fatalf("fecha de Lima incorrecta: %s", feedback.ServiceDate)
	}
}

func TestVehicleFeedbackSubmitter_TokenNoDisponible(t *testing.T) {
	qrRepo := &fakeQRRepo{resolveErr: domain.ErrQRNotAvailable}
	feedbackRepo := &fakeFeedbackRepo{}
	submitter := newSubmitter(qrRepo, feedbackRepo, nil)

	if _, err := submitter.Submit(context.Background(), "Abcdefghijklmnopqrstuv", fbInput()); !errors.Is(err, domain.ErrQRNotAvailable) {
		t.Fatalf("esperaba ErrQRNotAvailable, obtuve %v", err)
	}
	if feedbackRepo.createCalls != 0 {
		t.Fatal("no debe escribir cuando el token no está disponible")
	}
}

func TestVehicleFeedbackSubmitter_TokenConFormatoInvalido(t *testing.T) {
	feedbackRepo := &fakeFeedbackRepo{}
	submitter := newSubmitter(&fakeQRRepo{target: fbTarget()}, feedbackRepo, nil)

	if _, err := submitter.Submit(context.Background(), "invalido", fbInput()); !errors.Is(err, domain.ErrQRNotAvailable) {
		t.Fatalf("esperaba ErrQRNotAvailable, obtuve %v", err)
	}
	if feedbackRepo.createCalls != 0 {
		t.Fatal("no debe escribir con un token inválido")
	}
}

func TestVehicleFeedbackSubmitter_ValidacionesAntesDeEscribir(t *testing.T) {
	cases := []struct {
		name  string
		input VehicleFeedbackInput
		want  error
	}{
		{"rating bajo", VehicleFeedbackInput{Rating: 0, DeviceID: fbInput().DeviceID}, domain.ErrInvalidRating},
		{"rating alto", VehicleFeedbackInput{Rating: 6, DeviceID: fbInput().DeviceID}, domain.ErrInvalidRating},
		{"device inválido", VehicleFeedbackInput{Rating: 4, DeviceID: "xyz"}, domain.ErrInvalidDeviceID},
		{
			"comentario largo",
			VehicleFeedbackInput{Rating: 4, Comment: longComment(501), DeviceID: fbInput().DeviceID},
			domain.ErrInvalidFeedbackComment,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			feedbackRepo := &fakeFeedbackRepo{}
			submitter := newSubmitter(&fakeQRRepo{target: fbTarget()}, feedbackRepo, nil)
			if _, err := submitter.Submit(context.Background(), "Abcdefghijklmnopqrstuv", tc.input); !errors.Is(err, tc.want) {
				t.Fatalf("esperaba %v, obtuve %v", tc.want, err)
			}
			if feedbackRepo.createCalls != 0 {
				t.Fatal("las validaciones deben fallar antes de escribir")
			}
		})
	}
}

func TestVehicleFeedbackSubmitter_CalificacionDuplicada(t *testing.T) {
	feedbackRepo := &fakeFeedbackRepo{createErr: domain.ErrFeedbackAlreadySubmitted}
	submitter := newSubmitter(&fakeQRRepo{target: fbTarget()}, feedbackRepo, nil)

	if _, err := submitter.Submit(context.Background(), "Abcdefghijklmnopqrstuv", fbInput()); !errors.Is(err, domain.ErrFeedbackAlreadySubmitted) {
		t.Fatalf("esperaba ErrFeedbackAlreadySubmitted, obtuve %v", err)
	}
}

func TestVehicleFeedbackSubmitter_BusLiveConPosicion(t *testing.T) {
	positions := &fakePositionStore{position: &domain.Position{BusID: "bus-01"}}
	feedbackRepo := &fakeFeedbackRepo{}
	submitter := newSubmitter(&fakeQRRepo{target: fbTarget()}, feedbackRepo, positions)

	feedback, err := submitter.Submit(context.Background(), "Abcdefghijklmnopqrstuv", fbInput())
	if err != nil {
		t.Fatalf("envío falló: %v", err)
	}
	if !feedback.BusLive {
		t.Fatal("con posición reciente BusLive debe ser true")
	}
}

func TestVehicleFeedbackSubmitter_SinPosicionOSinRedisNoImpideElEnvio(t *testing.T) {
	cases := []struct {
		name      string
		positions ports.PositionStore
	}{
		{"sin posición", &fakePositionStore{}},
		{"redis caído", &fakePositionStore{err: errors.New("conexión rechazada")}},
		{"sin store", nil},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			feedbackRepo := &fakeFeedbackRepo{}
			submitter := newSubmitter(&fakeQRRepo{target: fbTarget()}, feedbackRepo, tc.positions)

			feedback, err := submitter.Submit(context.Background(), "Abcdefghijklmnopqrstuv", fbInput())
			if err != nil {
				t.Fatalf("el envío debe continuar: %v", err)
			}
			if feedback.BusLive {
				t.Fatal("sin posición BusLive debe ser false")
			}
		})
	}
}

func longComment(n int) string {
	runes := make([]rune, n)
	for i := range runes {
		runes[i] = 'á'
	}
	return string(runes)
}
