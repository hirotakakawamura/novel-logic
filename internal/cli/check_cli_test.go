package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"novel-logic/internal/lean"
)

func TestCheckQuick(t *testing.T) {
	dir := writeCLIProject(t)
	out, code := runCLI(t, "-C", dir, "check", "--quick")
	if code != 0 {
		t.Fatalf("exit code = %d, output = %q", code, out)
	}
	if !strings.Contains(out, "OK: stage1") {
		t.Fatalf("output = %q", out)
	}
}

func TestCheckExitsWhenLeanMissing(t *testing.T) {
	dir := writeCLIProject(t)
	oldPath := os.Getenv("PATH")
	t.Setenv("PATH", "")
	t.Cleanup(func() { _ = os.Setenv("PATH", oldPath) })

	out, code := runCLI(t, "-C", dir, "check")
	if code != 5 {
		t.Fatalf("exit code = %d, want 5; output = %q", code, out)
	}
	if !strings.Contains(out, "OK: stage1") {
		t.Fatalf("expected stage1 OK before lean skip, output = %q", out)
	}
}

func TestCheckStage2FailsOnBrokenLean(t *testing.T) {
	tc := lean.Detect()
	if !tc.Found {
		t.Skip("lean/lake not installed")
	}

	dir := writeCLIProject(t)
	if _, code := runCLI(t, "-C", dir, "generate"); code != 0 {
		t.Fatal("generate failed")
	}
	theorems := filepath.Join(dir, "logic", "Theorems.lean")
	if err := os.WriteFile(theorems, []byte("syntax error here\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	_, code := runCLI(t, "-C", dir, "check", "--no-generate")
	if code != 3 {
		t.Fatalf("exit %d, want 3 for stage2 failure", code)
	}
}

func TestDoctorMissingLean(t *testing.T) {
	dir := writeCLIProject(t)
	oldPath := os.Getenv("PATH")
	t.Setenv("PATH", "")
	t.Cleanup(func() { _ = os.Setenv("PATH", oldPath) })

	_, code := runCLI(t, "-C", dir, "doctor")
	if code != 5 {
		t.Fatalf("exit %d, want 5", code)
	}
}
