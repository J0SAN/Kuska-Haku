package main

import (
	"strings"
	"testing"
)

func TestCleanConnectionURL_QuitaSaltosDeLineaYEspacios(t *testing.T) {
	raw := " postgresql://usuario:clave@host.example.com/db?sslmode=r\n  equire \n"
	want := "postgresql://usuario:clave@host.example.com/db?sslmode=require"
	if got := cleanConnectionURL(raw); got != want {
		t.Fatalf("esperaba %q, obtuve %q", want, got)
	}
}

func TestRedactSecrets_OcultaLaClaveDeLaURL(t *testing.T) {
	msg := `parse "postgresql://usuario:Cl4ve-Secreta@host.example.com/db?sslmode=r\n": invalid control character`
	got := redactSecrets(msg)
	if strings.Contains(got, "Cl4ve-Secreta") {
		t.Fatalf("la clave sigue visible: %s", got)
	}
	if !strings.Contains(got, "postgresql://usuario:****@host.example.com") {
		t.Fatalf("debe conservar usuario y host para diagnosticar: %s", got)
	}
}

func TestRedactSecrets_OcultaPasswordEnFormatoClaveValor(t *testing.T) {
	got := redactSecrets("host=db user=app password=Cl4ve-Secreta sslmode=require")
	if strings.Contains(got, "Cl4ve-Secreta") || !strings.Contains(got, "password=****") {
		t.Fatalf("la clave debe ocultarse: %s", got)
	}
}

func TestRedactSecrets_SinClaveNoCambia(t *testing.T) {
	msg := "dial tcp 127.0.0.1:5432: connect: connection refused"
	if got := redactSecrets(msg); got != msg {
		t.Fatalf("no debía cambiar: %s", got)
	}
}
