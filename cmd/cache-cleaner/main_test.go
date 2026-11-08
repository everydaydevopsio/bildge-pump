package main

import (
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
