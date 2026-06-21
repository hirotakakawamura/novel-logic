package project

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestPinNovelRevisionAndHints(t *testing.T) {
	d := newTestProject(t)
	dir := d.Root
	gitInitTestRepo(t, dir)
	if err := d.AddNovel("scene1", MainBranch, "", true); err != nil {
		t.Fatal(err)
	}
	gitCommitAllTest(t, dir, "add novel")

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
	gitInitTestRepo(t, dir)
	if err := d.AddNovel("scene1", MainBranch, "", true); err != nil {
		t.Fatal(err)
	}
	gitCommitAllTest(t, dir, "v1")
	if _, err := d.PinNovelRevision("scene1", MainBranch, "", "", false); err != nil {
		t.Fatal(err)
	}

	body := filepath.Join(dir, DefaultNovelBodyPath("scene1", MainBranch))
	if err := os.WriteFile(body, []byte("edited\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	issues := NovelRevisionIssues(d)
	if !issueContains(issues, "uncommitted changes") {
		t.Fatalf("expected dirty issue, got %v", issues)
	}

	runGitTest(t, dir, "add", body)
	runGitTest(t, dir, "commit", "-m", "v2")
	issues = NovelRevisionIssues(d)
	if !issueContains(issues, "pinned revision") || !issueContains(issues, "differs from latest git commit") {
		t.Fatalf("expected drift issue, got %v", issues)
	}
}

func issueContains(issues []string, sub string) bool {
	for _, iss := range issues {
		if strings.Contains(iss, sub) {
			return true
		}
	}
	return false
}

func TestPinNovelRevisionRejectsDirty(t *testing.T) {
	d := newTestProject(t)
	dir := d.Root
	gitInitTestRepo(t, dir)
	if err := d.AddNovel("scene1", MainBranch, "", true); err != nil {
		t.Fatal(err)
	}
	gitCommitAllTest(t, dir, "v1")

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