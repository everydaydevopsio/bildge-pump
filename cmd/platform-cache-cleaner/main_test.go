package main

import (
	"path/filepath"
	"testing"
)

func TestLoadRequiresExplicitConfigWhenRequested(t *testing.T) {
	missing := filepath.Join(t.TempDir(), "missing.yaml")
	if _, err := load(missing, false); err == nil {
		t.Fatal("expected missing explicit config to fail")
	}
	if cfg, err := load(missing, true); err != nil || len(cfg.Targets) == 0 {
		t.Fatalf("implicit default config = %#v, %v; want default targets", cfg, err)
	}
}
