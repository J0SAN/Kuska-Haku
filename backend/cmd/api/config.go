package main

import "strconv"

// parseTrustedProxyHops interpreta TRUSTED_PROXY_HOPS. Un valor ausente o
// inválido equivale a 0: se usa la IP del socket y se ignora X-Forwarded-For.
func parseTrustedProxyHops(raw string) int {
	hops, err := strconv.Atoi(raw)
	if err != nil || hops < 0 {
		return 0
	}
	return hops
}
