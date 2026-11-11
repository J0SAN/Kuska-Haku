package main

import (
	"regexp"
	"strings"
)

// cleanConnectionURL elimina todo espacio en blanco de una URL de conexión.
// Una URL válida nunca lo contiene, pero al pegarla en un panel web es fácil
// que se cuele un salto de línea, y la API no arrancaría.
func cleanConnectionURL(raw string) string {
	return strings.Join(strings.Fields(raw), "")
}

var (
	// postgres://usuario:clave@host → postgres://usuario:****@host
	urlPasswordPattern = regexp.MustCompile(`(://[^:/@\s]*:)[^@\s]*@`)
	// password=clave (formato clave=valor de libpq)
	keywordPasswordPattern = regexp.MustCompile(`(?i)(password=)\S+`)
)

// redactSecrets oculta las contraseñas de las URLs de conexión antes de
// escribir un error en los logs: algunos drivers incluyen la URL completa en
// el mensaje de error.
func redactSecrets(msg string) string {
	msg = urlPasswordPattern.ReplaceAllString(msg, "${1}****@")
	return keywordPasswordPattern.ReplaceAllString(msg, "${1}****")
}
