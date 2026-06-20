package project

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultNovelBodyPath(t *testing.T) {
	got := DefaultNovelBodyPath("scene1", MainBranch)
	want := filepath.Join(DirNovels, MainBranch, "scene1"+NovelBodyExt)
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestAddNovelPerBranch(t *testing.T) {
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
	body := filepath.Join(d.Root, DefaultNovelBodyPath("scene1", "branch_a"))
	if _, err := os.Stat(body); err != nil {
		t.Fatalf("expected body file: %v", err)
	}
	if err := d.AddNovel("scene1", MainBranch, "", false); err == nil {
		t.Fatal("expected duplicate novel on same branch")
	}
}

func TestFindNovelByBranch(t *testing.T) {
	d := newTestProject(t)
	if err := d.AddNovel("scene1", MainBranch, "", true); err != nil {
		t.Fatal(err)
	}
	n, idx := d.FindNovel("scene1", MainBranch)
	if n == nil || idx < 0 {
		t.Fatal("novel not found")
	}
	if _, idx := d.FindNovel("scene1", "branch_a"); idx >= 0 {
		t.Fatal("novel should not exist on branch_a")
	}
}