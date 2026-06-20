package generate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"novel-logic/internal/project"
)

var snapshotFiles = []string{
	"Project.lean",
	"Facts.lean",
	"Rules.lean",
	"Timeline.lean",
	"Theorems.lean",
	"lakefile.toml",
}

func TestRunMinimalSnapshot(t *testing.T) {
	d := minimalProject(t)
	if err := Run(d); err != nil {
		t.Fatal(err)
	}

	for _, name := range snapshotFiles {
		gotPath := filepath.Join(d.Root, project.DirLogic, name)
		got, err := os.ReadFile(gotPath)
		if err != nil {
			t.Fatalf("read generated %s: %v", name, err)
		}

		goldenPath := filepath.Join("testdata", "minimal", name)
		if os.Getenv("UPDATE_GOLDEN") != "" {
			if err := os.MkdirAll(filepath.Dir(goldenPath), 0o755); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(goldenPath, got, 0o644); err != nil {
				t.Fatal(err)
			}
			continue
		}

		want, err := os.ReadFile(goldenPath)
		if err != nil {
			t.Fatalf("read golden %s: %v (run UPDATE_GOLDEN=1 go test ./internal/generate)", goldenPath, err)
		}
		if string(got) != string(want) {
			t.Fatalf("snapshot mismatch for %s", name)
		}
	}
}

func TestRunCreatesCoreLean(t *testing.T) {
	d := minimalProject(t)
	if err := Run(d); err != nil {
		t.Fatal(err)
	}
	corePath := filepath.Join(d.Root, project.DirLogic, "Core.lean")
	b, err := os.ReadFile(corePath)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(b), "namespace NovelLogic") {
		t.Fatal("Core.lean missing NovelLogic namespace")
	}
}

func TestNamespaceASCII(t *testing.T) {
	d := minimalProject(t)
	if got := namespace(d); got != "test" {
		t.Fatalf("namespace = %q, want test", got)
	}
}

func TestNamespaceJapanese(t *testing.T) {
	d := minimalProject(t)
	d.Meta.Title = "桃太郎"
	if got := namespace(d); got != "Momotaro" {
		t.Fatalf("namespace = %q, want Momotaro", got)
	}
}