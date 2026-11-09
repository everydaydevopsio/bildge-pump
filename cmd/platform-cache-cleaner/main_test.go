package main

import (
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
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

func TestWindowsBatchLauncherPreservesExecutableAndArgumentBoundaries(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("native Windows batch invocation test")
	}
	executable := filepath.Join(t.TempDir(), "tool with spaces.cmd")
	if err := os.WriteFile(executable, []byte("@echo off\r\necho %*\r\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	args := []string{"argument with spaces", "literal & echo INJECTED", "pipe|caret^", "100%TEMP%"}
	cmd, err := windowsBatchCommand(executable, args)
	if err != nil {
		t.Fatal(err)
	}
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("batch command failed: %v: %s", err, output)
	}
	text := strings.TrimSpace(string(output))
	if strings.Count(text, "\n") != 0 || !strings.Contains(text, `"argument with spaces"`) || !strings.Contains(text, `"literal & echo INJECTED"`) || !strings.Contains(text, `"pipe|caret^"`) || !strings.Contains(text, `"100%TEMP%"`) {
		t.Fatalf("arguments were changed or split: %q", text)
	}
	if _, err := windowsBatchCommand(executable, []string{`literal " quote`}); err == nil {
		t.Fatal("expected unrepresentable literal quote to be rejected")
	}
}
