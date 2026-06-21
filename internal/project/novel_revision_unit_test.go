package project

import (
	"os"
	"path/filepath"
	"testing"
)

func TestPinNovelRevisionAndHints(t *testing.T) {
	d := newTestProject(t)
	dir := d.Root
	runGitTest(t, dir, "init")
	runGitTest(t, dir, "config", "user.email", "test@example.com")
	runGitTest(t, dir, "config", "user.name", "Test User")
	if err := d.AddNovel("scene1", MainBranch, "", true); err != nil {
		t.Fatal(err)
	}
	runGitTest(t, dir, "add", ".")
	runGitTest(t, dir, "commit", "-m", "add novel")

	if hints := NovelRevisionHints(d); len(hints) == 0 {
		t.Fatal("expected unpinned hint")
	}

	entry, err := d.PinNovelRevision("scene1", MainBranch, "", "first pin", false)
	if err != nil {
		t.Fatal(err)
	}
	if entry.Revision == "" || entry.Short == "" {
		t.Fatalf("entry = %+v", entry)
	}
	if hints := NovelRevisionHints(d); len(hints) != 0 {
		t.Fatalf("hints = %v", hints)
	}
}

func TestNovelRevisionIssuesOnDriftAndDirty(t *testing.T) {
	d := newTestProject(t)
	dir := d.Root
	runGitTest(t, dir, "init")
	runGitTest(t, dir, "config", "user.email", "test@example.com")
	runGitTest(t, dir, "config", "user.name", "Test User")
	if err := d.AddNovel("scene1", MainBranch, "", true); err != nil {
		t.Fatal(err)
	}
	runGitTest(t, dir, "add", ".")
	runGitTest(t, dir, "commit", "-m", "v1")
	if _, err := d.PinNovelRevision("scene1", MainBranch, "", "", false); err != nil {
		t.Fatal(err)
	}

	body := filepath.Join(dir, DefaultNovelBodyPath("scene1", MainBranch))
	if err := os.WriteFile(body, []byte("edited\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if issues := NovelRevisionIssues(d); len(issues) == 0 {
		t.Fatal("expected dirty or drift issue")
	}

	runGitTest(t, dir, "add", body)
	runGitTest(t, dir, "commit", "-m", "v2")
	if issues := NovelRevisionIssues(d); len(issues) == 0 {
		t.Fatal("expected revision drift issue")
	}
}

func TestPinNovelRevisionRejectsDirty(t *testing.T) {
	d := newTestProject(t)
	dir := d.Root
	runGitTest(t, dir, "init")
	runGitTest(t, dir, "config", "user.email", "test@example.com")
	runGitTest(t, dir, "config", "user.name", "Test User")
	if err := d.AddNovel("scene1", MainBranch, "", true); err != nil {
		t.Fatal(err)
	}
	runGitTest(t, dir, "add", ".")
	runGitTest(t, dir, "commit", "-m", "v1")

	body := filepath.Join(dir, DefaultNovelBodyPath("scene1", MainBranch))
	if err := os.WriteFile(body, []byte("dirty edit\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := d.PinNovelRevision("scene1", MainBranch, "", "", false); err == nil {
		t.Fatal("expected dirty error")
	}
	if _, err := d.PinNovelRevision("scene1", MainBranch, "", "", true); err != nil {
		t.Fatal(err)
	}
}