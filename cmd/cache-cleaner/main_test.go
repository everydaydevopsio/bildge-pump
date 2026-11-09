package main

import (
	"bytes"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestPlanTaskArgsRoutesYesOnlyToDevCache(t *testing.T) {
	got, err := planTaskArgs([]string{"dev-cache", "git-cleaner", "linux-cache-cleaner"}, nil, true, "/tmp/src", true)
	if err != nil {
		t.Fatal(err)
	}
	want := map[string][]string{
		"dev-cache":           {"--scan", "/tmp/src", "--clean", "--yes"},
		"git-cleaner":         {"--scan", "/tmp/src", "--clean"},
		"linux-cache-cleaner": {"--clean"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("planTaskArgs() = %#v, want %#v", got, want)
	}
}

func TestPlanTaskArgsRejectsUnsupportedArgumentBeforeLaunching(t *testing.T) {
	_, err := planTaskArgs([]string{"dev-cache", "git-cleaner"}, []string{"--depth", "2"}, true, "/tmp/src", false)
	if err == nil {
		t.Fatal("expected --depth to be rejected for a multi-task cleanup")
	}
}

func TestParseBooleanFlagHonorsExplicitValue(t *testing.T) {
	for _, tt := range []struct {
		name     string
		value    string
		hasValue bool
		want     bool
		wantErr  bool
	}{
		{name: "implicit true", want: true},
		{name: "explicit true", value: "true", hasValue: true, want: true},
		{name: "explicit false", value: "false", hasValue: true, want: false},
		{name: "invalid", value: "sometimes", hasValue: true, wantErr: true},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseBooleanFlag(tt.value, tt.hasValue)
			if (err != nil) != tt.wantErr || got != tt.want {
				t.Fatalf("parseBooleanFlag() = %v, %v; want %v, error=%v", got, err, tt.want, tt.wantErr)
			}
		})
	}
}

func TestLoadConfigWithDefaultsRejectsMissingExplicitConfig(t *testing.T) {
	missing := filepath.Join(t.TempDir(), "missing.yaml")
	if _, err := loadConfigWithDefaults(missing, false); err == nil {
		t.Fatal("expected missing explicit config to fail")
	}
	if cfg, err := loadConfigWithDefaults(missing, true); err != nil || len(cfg.Tasks) == 0 {
		t.Fatalf("implicit default config = %#v, %v; want default tasks", cfg, err)
	}
}

func TestRunMissingExplicitConfigStopsBeforeTasks(t *testing.T) {
	missing := filepath.Join(t.TempDir(), "missing.yaml")
	oldStderr := os.Stderr
	reader, writer, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	os.Stderr = writer
	code := run([]string{"--config", missing, "clean", "--all"})
	os.Stderr = oldStderr
	_ = writer.Close()
	var output bytes.Buffer
	_, _ = output.ReadFrom(reader)
	_ = reader.Close()
	if code == 0 {
		t.Fatal("run succeeded with a missing explicit config")
	}
	if !bytes.Contains(output.Bytes(), []byte("config file "+missing+" does not exist")) {
		t.Fatalf("expected config error before launching tasks, got %q", output.String())
	}
}

func TestParseBooleanFalseDoesNotRequestCleanup(t *testing.T) {
	clean, err := parseBooleanFlag("false", true)
	if err != nil {
		t.Fatal(err)
	}
	if clean {
		t.Fatal("--clean=false must not request cleanup")
	}
	yes, err := parseBooleanFlag("false", true)
	if err != nil {
		t.Fatal(err)
	}
	if yes {
		t.Fatal("--yes=false must not skip confirmation")
	}
	if cleanupRequested("scan", true, clean, nil) {
		t.Fatal("scan --clean=false must remain a dry run")
	}
	if cleanupRequested("clean", true, clean, nil) {
		t.Fatal("clean --clean=false must remain a dry run")
	}
	if !cleanupRequested("clean", false, false, nil) {
		t.Fatal("clean command should request cleanup without a --clean override")
	}
}
