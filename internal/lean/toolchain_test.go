package lean

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDetectNotFoundWithEmptyPath(t *testing.T) {
	oldPath := os.Getenv("PATH")
	t.Setenv("PATH", "")
	t.Cleanup(func() { _ = os.Setenv("PATH", oldPath) })

	tc := Detect()
	if tc.Found {
		t.Fatal("expected Found=false with empty PATH")
	}
	if tc.Version() != "not found" {
		t.Fatalf("Version() = %q", tc.Version())
	}
}

func TestLakeBuildNotFound(t *testing.T) {
	oldPath := os.Getenv("PATH")
	t.Setenv("PATH", "")
	t.Cleanup(func() { _ = os.Setenv("PATH", oldPath) })

	dir := t.TempDir()
	out, err := LakeBuild(dir, 0)
	if err == nil {
		t.Fatal("expected error when lean/lake missing")
	}
	if !strings.Contains(err.Error(), "not found") {
		t.Fatalf("err = %v", err)
	}
	if out != "" {
		t.Fatalf("expected empty output, got %q", out)
	}
}

func TestDetectFindsToolsWhenPresent(t *testing.T) {
	tc := Detect()
	if !tc.Found {
		t.Skip("lean/lake not installed; skipping presence test")
	}
	if tc.Lean == "" || tc.Lake == "" {
		t.Fatalf("tc = %+v", tc)
	}
	if !strings.Contains(tc.Version(), "Lean") && tc.Version() != "not found" {
		// lean --version usually contains "Lean"
		if tc.Version() == "" {
			t.Fatalf("empty version string")
		}
	}
	_ = filepath.Dir(tc.Lake)
}