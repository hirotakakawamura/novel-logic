package project

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNovelBodyIssuesMissingFile(t *testing.T) {
	d := newTestProject(t)
	body := DefaultNovelBodyPath("scene1", MainBranch)
	d.Novels = append(d.Novels, NovelMeta{
		SceneID: "scene1", Branch: MainBranch,
		TimeStart: "t1", TimeEnd: "t2", BodyPath: body,
	})
	issues := NovelBodyIssues(d)
	if len(issues) == 0 {
		t.Fatal("expected missing body issue")
	}
}

func TestUpdateNovelSyncsSceneWindow(t *testing.T) {
	d := newTestProject(t)
	if err := d.AddNovel("scene1", MainBranch, "", true); err != nil {
		t.Fatal(err)
	}
	for i := range d.Scenes {
		if d.Scenes[i].ID == "scene1" {
			d.Scenes[i].TimeEnd = "t3"
		}
	}
	if err := d.UpdateNovel("scene1", MainBranch, "", false); err != nil {
		t.Fatal(err)
	}
	n, _ := d.FindNovel("scene1", MainBranch)
	if n.TimeEnd != "t3" {
		t.Fatalf("TimeEnd = %q, want t3", n.TimeEnd)
	}
}

func TestRemoveNovelDeletesBodyWhenRequested(t *testing.T) {
	d := newTestProject(t)
	if err := d.AddNovel("scene1", MainBranch, "", true); err != nil {
		t.Fatal(err)
	}
	body := filepath.Join(d.Root, DefaultNovelBodyPath("scene1", MainBranch))
	if err := d.RemoveNovel("scene1", MainBranch, false); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(body); !os.IsNotExist(err) {
		t.Fatalf("body file should be deleted: %v", err)
	}
}