package main

import "testing"

func TestParseFlagSupportsSeparateAndEqualsForms(t *testing.T) {
	args := []string{"--config", "test.toml", "--chat=12345"}

	if got := parseFlag(args, "--config", "default.toml"); got != "test.toml" {
		t.Fatalf("--config ayrik formu beklenmedi: got=%q", got)
	}
	if got := parseFlag(args, "--chat", "0"); got != "12345" {
		t.Fatalf("--chat equals formu beklenmedi: got=%q", got)
	}
	if got := parseFlag(args, "--missing", "fallback"); got != "fallback" {
		t.Fatalf("eksik bayrakta varsayilan donmeli: got=%q", got)
	}
}

func TestParseFlagMissingValueFallsBackToDefault(t *testing.T) {
	args := []string{"--config"}

	if got := parseFlag(args, "--config", "default.toml"); got != "default.toml" {
		t.Fatalf("degersiz bayrakta varsayilan donmeli: got=%q", got)
	}
}
