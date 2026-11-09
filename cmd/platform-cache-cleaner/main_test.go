package main

import (
	"path/filepath"
	"reflect"
	"testing"
)

func TestScanTargetReportIncludesPathlessDockerCommands(t *testing.T) {
	target := Target{Name: "docker", Enabled: true, Cmds: [][]string{{"docker", "builder", "prune", "-af"}, {"docker", "system", "prune", "-af"}}}
	report := Report{Findings: map[string][]Finding{}, Totals: map[string]int64{}}
	addTargetToReport(&report, target)
	if findings, ok := report.Findings["docker"]; !ok || len(findings) != 0 {
		t.Fatalf("docker findings = %#v, present=%v; want empty entry", findings, ok)
	}
	if !reflect.DeepEqual(report.PlannedCommands["docker"], target.Cmds) {
		t.Fatalf("planned Docker commands = %#v, want %#v", report.PlannedCommands["docker"], target.Cmds)
	}
}

func TestLoadRequiresExplicitConfigWhenRequested(t *testing.T) {
	missing := filepath.Join(t.TempDir(), "missing.yaml")
	if _, err := load(missing, false); err == nil {
		t.Fatal("expected missing explicit config to fail")
	}
	if cfg, err := load(missing, true); err != nil || len(cfg.Targets) == 0 {
		t.Fatalf("implicit default config = %#v, %v; want default targets", cfg, err)
	}
}
