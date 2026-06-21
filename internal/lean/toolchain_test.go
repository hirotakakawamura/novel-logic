package lean

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"novel-logic/internal/generate"
	"novel-logic/internal/project"
	"novel-logic/internal/testfixture"
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

func TestLakeBuildSuccess(t *testing.T) {
	tc := Detect()
	if !tc.Found {
		t.Skip("lean/lake not installed")
	}

	logicDir := walkthroughLogicDir(t)
	out, err := LakeBuild(logicDir, 0)
	if err != nil {
		t.Fatalf("LakeBuild: %v\noutput:\n%s", err, out)
	}
}

func TestLakeBuildPassesJobsFlag(t *testing.T) {
	tmpBin := t.TempDir()
	logPath := filepath.Join(t.TempDir(), "lake-args.log")
	fakeLean := filepath.Join(tmpBin, "lean")
	fakeLake := filepath.Join(tmpBin, "lake")
	if err := os.WriteFile(fakeLean, []byte("#!/bin/sh\nexit 0\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	lakeScript := "#!/bin/sh\necho \"$@\" >> " + logPath + "\nexit 0\n"
	if err := os.WriteFile(fakeLake, []byte(lakeScript), 0o755); err != nil {
		t.Fatal(err)
	}

	oldPath := os.Getenv("PATH")
	t.Setenv("PATH", tmpBin+string(os.PathListSeparator)+oldPath)
	t.Cleanup(func() { _ = os.Setenv("PATH", oldPath) })

	logicDir := t.TempDir()
	out, err := LakeBuild(logicDir, 2)
	if err != nil {
		t.Fatalf("LakeBuild: %v\noutput:\n%s", err, out)
	}
	args, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatal(err)
	}
	got := string(args)
	if !strings.Contains(got, "-j") || !strings.Contains(got, "2") {
		t.Fatalf("lake args = %q, want -j 2", got)
	}
}

func TestLakeBuildFailsOnInvalidLean(t *testing.T) {
	tc := Detect()
	if !tc.Found {
		t.Skip("lean/lake not installed")
	}

	logicDir := generatedLogicDir(t)
	theorems := filepath.Join(logicDir, "Theorems.lean")
	if err := os.WriteFile(theorems, []byte("syntax error here\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	out, err := LakeBuild(logicDir, 0)
	if err == nil {
		t.Fatalf("expected build failure, output:\n%s", out)
	}
	if strings.TrimSpace(out) == "" {
		t.Fatal("expected non-empty build output on failure")
	}
}

func TestVersionReturnsOutputWhenLeanFails(t *testing.T) {
	tmpBin := t.TempDir()
	fakeLean := filepath.Join(tmpBin, "lean")
	script := "#!/bin/sh\necho broken lean version\nexit 1\n"
	if err := os.WriteFile(fakeLean, []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}

	oldPath := os.Getenv("PATH")
	t.Setenv("PATH", tmpBin+string(os.PathListSeparator)+oldPath)
	t.Cleanup(func() { _ = os.Setenv("PATH", oldPath) })

	tc := Detect()
	if tc.Lean != fakeLean {
		t.Fatalf("Detect().Lean = %q, want fake %q", tc.Lean, fakeLean)
	}
	if got := tc.Version(); got != "broken lean version" {
		t.Fatalf("Version() = %q, want broken lean version", got)
	}
}

func walkthroughLogicDir(t *testing.T) string {
	t.Helper()
	root := moduleRoot(t)
	logicDir := filepath.Join(root, "examples", "momotaro-walkthrough", "logic")
	if _, err := os.Stat(filepath.Join(logicDir, "lakefile.toml")); err != nil {
		t.Fatalf("walkthrough logic dir: %v", err)
	}
	return logicDir
}

func generatedLogicDir(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	testfixture.WriteFiles(t, dir, testfixture.MinimalFiles())
	d, err := project.Load(dir)
	if err != nil {
		t.Fatal(err)
	}
	if err := generate.Run(d); err != nil {
		t.Fatal(err)
	}
	return filepath.Join(dir, project.DirLogic)
}

func moduleRoot(t *testing.T) string {
	t.Helper()
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatal("go.mod not found")
		}
		dir = parent
	}
}