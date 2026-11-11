package main

import "testing"

func TestParseTrustedProxyHops(t *testing.T) {
	cases := map[string]int{
		"":    0,
		"0":   0,
		"1":   1,
		"2":   2,
		"-1":  0,
		"abc": 0,
	}
	for raw, want := range cases {
		if got := parseTrustedProxyHops(raw); got != want {
			t.Fatalf("parseTrustedProxyHops(%q) = %d, esperaba %d", raw, got, want)
		}
	}
}
