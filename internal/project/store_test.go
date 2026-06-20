package project

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSaveLoadRoundtrip(t *testing.T) {
	d := newTestProject(t)
	if err := d.AddBranch("branch_a", "alt", MainBranch, "", ""); err != nil {
		t.Fatal(err)
	}
	if err := d.AddNovel("scene1", MainBranch, "", true); err != nil {
		t.Fatal(err)
	}
	if err := d.AddNovel("scene1", "branch_a", "", true); err != nil {
		t.Fatal(err)
	}
	_, err := d.AddAction("hero", "mid", "alt_end", "t3", "plot", "", "branch_a")
	if err != nil {
		t.Fatal(err)
	}
	d.Facts = append(d.Facts, Fact{
		ID: "fact_branch", Kind: FactState, Thing: "hero", Pred: "alt_only", Scope: "plot", Branch: "branch_a",
	})

	if err := Save(d); err != nil {
		t.Fatal(err)
	}

	loaded, err := Load(d.Root)
	if err != nil {
		t.Fatal(err)
	}

	if loaded.Meta.Title != d.Meta.Title {
		t.Fatalf("title: got %q, want %q", loaded.Meta.Title, d.Meta.Title)
	}
	if len(loaded.Branches) != len(d.Branches) {
		t.Fatalf("branches: got %d, want %d", len(loaded.Branches), len(d.Branches))
	}
	if _, idx := loaded.FindBranchDef("branch_a"); idx < 0 {
		t.Fatal("branch_a not found after reload")
	}
	if _, idx := loaded.FindBranchDef(MainBranch); idx < 0 {
		t.Fatal("main branch not found after reload")
	}
	if len(loaded.Novels) != 2 {
		t.Fatalf("novels: got %d, want 2", len(loaded.Novels))
	}
	for _, n := range loaded.Novels {
		want := DefaultNovelBodyPath(n.SceneID, n.Branch)
		if n.BodyPath != want {
			t.Fatalf("novel %s/%s body_path: got %q, want %q", n.SceneID, n.Branch, n.BodyPath, want)
		}
	}
	if len(loaded.Actions) != len(d.Actions) {
		t.Fatalf("actions: got %d, want %d", len(loaded.Actions), len(d.Actions))
	}
	for _, a := range loaded.Actions {
		if a.Branch == "" {
			t.Fatalf("action %q has empty branch after reload", a.ID)
		}
	}
	for _, f := range loaded.Facts {
		if f.Branch == "" && f.ID == "fact_branch" {
			t.Fatal("branch fact should retain branch_a")
		}
	}
}

func TestLoadEnsuresMainBranchOnEmptyBranchesFile(t *testing.T) {
	dir := t.TempDir()
	writeTestFile(t, dir, FileProject, "title: x\ntime_order: [t1]\n")
	writeTestFile(t, dir, FileThings, "[]\n")
	writeTestFile(t, dir, FileTimes, "- id: t1\n")
	writeTestFile(t, dir, FileScenes, "[]\n")
	writeTestFile(t, dir, FileBranches, "[]\n")
	writeTestFile(t, dir, FileForks, "[]\n")
	writeTestFile(t, dir, FileMerges, "[]\n")
	writeTestFile(t, dir, FileFacts, "[]\n")
	writeTestFile(t, dir, FileActions, "[]\n")
	writeTestFile(t, dir, FileRules, "[]\n")
	writeTestFile(t, dir, FileNovels, "[]\n")

	d, err := Load(dir)
	if err != nil {
		t.Fatal(err)
	}
	b, idx := d.FindBranchDef(MainBranch)
	if idx < 0 || b == nil {
		t.Fatal("main branch should be auto-added")
	}
	if b.Label != "本線" {
		t.Fatalf("main label = %q", b.Label)
	}
}

func TestLoadMigratesLegacyNovelBodyPath(t *testing.T) {
	dir := t.TempDir()
	writeTestFile(t, dir, FileProject, "title: x\ntime_order: [t1]\n")
	writeTestFile(t, dir, FileThings, "[]\n")
	writeTestFile(t, dir, FileTimes, "- id: t1\n")
	writeTestFile(t, dir, FileScenes, `- id: scene1
  summary: s
  time_start: t1
  time_end: t1
`)
	writeTestFile(t, dir, FileBranches, `- id: main
`)
	writeTestFile(t, dir, FileForks, "[]\n")
	writeTestFile(t, dir, FileMerges, "[]\n")
	writeTestFile(t, dir, FileFacts, "[]\n")
	writeTestFile(t, dir, FileActions, "[]\n")
	writeTestFile(t, dir, FileRules, "[]\n")
	writeTestFile(t, dir, FileNovels, `- scene_id: scene1
  time_start: t1
  time_end: t1
  body_path: novels/scene1.txt
`)
	legacyBody := filepath.Join(dir, "novels", "scene1.txt")
	if err := os.MkdirAll(filepath.Dir(legacyBody), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(legacyBody, []byte("body"), 0o644); err != nil {
		t.Fatal(err)
	}

	d, err := Load(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(d.Novels) != 1 {
		t.Fatalf("novels: got %d", len(d.Novels))
	}
	want := DefaultNovelBodyPath("scene1", MainBranch)
	if d.Novels[0].BodyPath != want {
		t.Fatalf("body_path: got %q, want %q", d.Novels[0].BodyPath, want)
	}
	if d.Novels[0].Branch != MainBranch {
		t.Fatalf("branch: got %q", d.Novels[0].Branch)
	}
}

func TestSavePersistsBranchFields(t *testing.T) {
	d := newTestProject(t)
	if err := d.AddBranch("branch_a", "alt", MainBranch, "", ""); err != nil {
		t.Fatal(err)
	}
	_, err := d.AddAction("hero", "mid", "alt_end", "t3", "plot", "", "branch_a")
	if err != nil {
		t.Fatal(err)
	}
	if err := Save(d); err != nil {
		t.Fatal(err)
	}

	raw, err := os.ReadFile(filepath.Join(d.Root, FileActions))
	if err != nil {
		t.Fatal(err)
	}
	content := string(raw)
	if !strings.Contains(content, "branch_a") || !strings.Contains(content, "alt_end") {
		t.Fatalf("actions.yaml missing branch data: %s", content)
	}
}