package validate

import (
	"testing"

	"novel-logic/internal/project"
	"novel-logic/internal/testfixture"
)

func TestHintsPlotSceneAlignment(t *testing.T) {
	d := testfixture.LoadMinimal(t)
	hints := Hints(d)
	if !hasIssueCode(hints, "action.plot_scene_hint") {
		t.Fatalf("expected action.plot_scene_hint, got %v", hints)
	}
}

func TestHintsNovelRevisionUnpinned(t *testing.T) {
	dir := testfixture.WriteMinimalDir(t)
	initGitRepo(t, dir)
	d, err := project.Load(dir)
	if err != nil {
		t.Fatal(err)
	}
	if err := d.AddNovel("scene1", project.MainBranch, "", true); err != nil {
		t.Fatal(err)
	}
	hints := Hints(d)
	if !hasIssueCode(hints, "novel.revision_hint") {
		t.Fatalf("expected novel.revision_hint, got %v", hints)
	}
}

func initGitRepo(t *testing.T, dir string) {
	t.Helper()
	runGit(t, dir, "init")
	runGit(t, dir, "config", "user.email", "test@example.com")
	runGit(t, dir, "config", "user.name", "Test User")
	runGit(t, dir, "add", ".")
	runGit(t, dir, "commit", "-m", "init")
}