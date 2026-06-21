package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestVersion(t *testing.T) {
	out, code := runCLI(t, "version")
	if code != 0 {
		t.Fatalf("exit code = %d", code)
	}
	if !strings.Contains(out, "novel-logic") {
		t.Fatalf("output = %q", out)
	}
}

func TestInfo(t *testing.T) {
	dir := writeCLIProject(t)
	out, code := runCLI(t, "-C", dir, "info")
	if code != 0 {
		t.Fatalf("exit code = %d, output = %q", code, out)
	}
	if !strings.Contains(out, "title: test") {
		t.Fatalf("output = %q", out)
	}
}

func TestPlotSet(t *testing.T) {
	dir := writeCLIProject(t)
	_, code := runCLI(t, "-C", dir, "plot", "set", "--title", "My Novel", "--summary", "A tale")
	if code != 0 {
		t.Fatalf("exit code = %d", code)
	}
	out, code := runCLI(t, "-C", dir, "info")
	if code != 0 || !strings.Contains(out, "title: My Novel") {
		t.Fatalf("info exit %d, output=%q", code, out)
	}
	out, code = runCLI(t, "-C", dir, "plot", "show")
	if code != 0 || !strings.Contains(out, "summary: A tale") {
		t.Fatalf("plot show exit %d, output=%q", code, out)
	}
}

func TestGenerate(t *testing.T) {
	dir := writeCLIProject(t)
	out, code := runCLI(t, "-C", dir, "generate")
	if code != 0 {
		t.Fatalf("exit code = %d, output = %q", code, out)
	}
	if !strings.Contains(out, "generated logic/") {
		t.Fatalf("output = %q", out)
	}
	for _, name := range []string{"Project.lean", "Facts.lean", "Theorems.lean", "lakefile.toml"} {
		path := filepath.Join(dir, "logic", name)
		if _, err := os.Stat(path); err != nil {
			t.Fatalf("missing %s: %v", name, err)
		}
	}
}
