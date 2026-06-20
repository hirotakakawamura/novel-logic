package template

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"novel-logic/internal/project"
)

func TestListIncludesBuiltinTemplates(t *testing.T) {
	names, err := List()
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"default", "momotaro"} {
		if !slices.Contains(names, want) {
			t.Fatalf("List() = %v, missing %q", names, want)
		}
	}
}

func TestMaterializeDefault(t *testing.T) {
	dir := t.TempDir()
	if err := Materialize(dir, "default"); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{
		project.FileProject, project.FileForks, project.FileMerges, project.FileBranches,
	} {
		if _, err := os.Stat(filepath.Join(dir, name)); err != nil {
			t.Fatalf("missing %s: %v", name, err)
		}
	}
}

func TestMaterializeMomotaro(t *testing.T) {
	dir := t.TempDir()
	if err := Materialize(dir, "momotaro"); err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile(filepath.Join(dir, project.FileProject))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(b), "桃太郎") {
		t.Fatalf("momotaro template project.yaml = %s", b)
	}
}

func TestMaterializeUnknownTemplate(t *testing.T) {
	dir := t.TempDir()
	if err := Materialize(dir, "no_such"); err == nil {
		t.Fatal("expected error for unknown template")
	}
}