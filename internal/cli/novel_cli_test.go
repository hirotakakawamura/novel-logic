package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"novel-logic/internal/project"
)

func TestNovelListEmpty(t *testing.T) {
	dir := writeCLIProject(t)
	out, code := runCLI(t, "-C", dir, "novel", "list")
	if code != 0 || !strings.Contains(out, "(none)") {
		t.Fatalf("exit %d, output=%q", code, out)
	}
}

func TestNovelListAndShow(t *testing.T) {
	dir := writeCLIProject(t)
	gitInit(t, dir)
	if _, code := runCLI(t, "-C", dir, "novel", "add", "scene1", "--init"); code != 0 {
		t.Fatal("novel add failed")
	}
	gitCommitAll(t, dir, "add novel")

	out, code := runCLI(t, "-C", dir, "novel", "list")
	if code != 0 || !strings.Contains(out, "scene1") {
		t.Fatalf("list exit %d, output=%q", code, out)
	}
	if !strings.Contains(out, "scene2 [no novel]") {
		t.Fatalf("output=%q", out)
	}

	out, code = runCLI(t, "-C", dir, "novel", "show", "scene1")
	if code != 0 || !strings.Contains(out, "body:") {
		t.Fatalf("show exit %d, output=%q", code, out)
	}
	if !strings.Contains(out, "alignment:") {
		t.Fatalf("output=%q", out)
	}
}

func TestNovelShowBodyTruncation(t *testing.T) {
	dir := writeCLIProject(t)
	if _, code := runCLI(t, "-C", dir, "novel", "add", "scene1", "--init"); code != 0 {
		t.Fatal("novel add failed")
	}
	body := filepath.Join(dir, "novels", "main", "scene1.txt")
	if err := os.WriteFile(body, []byte(strings.Repeat("x", 500)), 0o644); err != nil {
		t.Fatal(err)
	}
	out, code := runCLI(t, "-C", dir, "novel", "show", "scene1")
	if code != 0 || !strings.Contains(out, "…") {
		t.Fatalf("exit %d, output=%q", code, out)
	}
	out, code = runCLI(t, "-C", dir, "novel", "show", "scene1", "--full")
	if code != 0 || !strings.Contains(out, strings.Repeat("x", 100)) {
		t.Fatalf("full exit %d", code)
	}
}

func TestNovelUpdateSyncsSceneWindow(t *testing.T) {
	dir := writeCLIProject(t)
	if _, code := runCLI(t, "-C", dir, "novel", "add", "scene1", "--init"); code != 0 {
		t.Fatalf("novel add failed")
	}
	scenesPath := filepath.Join(dir, project.FileScenes)
	b, err := os.ReadFile(scenesPath)
	if err != nil {
		t.Fatal(err)
	}
	updated := strings.Replace(string(b), "time_end: t2", "time_end: t3", 1)
	if err := os.WriteFile(scenesPath, []byte(updated), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, code := runCLI(t, "-C", dir, "novel", "update", "scene1"); code != 0 {
		t.Fatalf("novel update exit %d", code)
	}
	d, err := project.Load(dir)
	if err != nil {
		t.Fatal(err)
	}
	n, _ := d.FindNovel("scene1", project.MainBranch)
	if n == nil || n.TimeEnd != "t3" {
		t.Fatalf("novel TimeEnd = %q, want t3", n.TimeEnd)
	}
}

func TestNovelRemoveKeepsBody(t *testing.T) {
	dir := writeCLIProject(t)
	if _, code := runCLI(t, "-C", dir, "novel", "add", "scene1", "--init"); code != 0 {
		t.Fatalf("novel add failed")
	}
	bodyAbs := filepath.Join(dir, "novels", "main", "scene1.txt")
	if _, code := runCLI(t, "-C", dir, "novel", "remove", "scene1", "--keep-body"); code != 0 {
		t.Fatalf("novel remove exit %d", code)
	}
	if _, err := os.Stat(bodyAbs); err != nil {
		t.Fatalf("body file should remain: %v", err)
	}
	d, err := project.Load(dir)
	if err != nil {
		t.Fatal(err)
	}
	if _, idx := d.FindNovel("scene1", project.MainBranch); idx >= 0 {
		t.Fatal("novel metadata should be removed")
	}
}

func TestNovelRevisionPinCLI(t *testing.T) {
	dir := writeCLIProject(t)
	gitInit(t, dir)
	if _, code := runCLI(t, "-C", dir, "novel", "add", "scene1", "--init"); code != 0 {
		t.Fatal("novel add failed")
	}
	gitCommitAll(t, dir, "add novel")

	out, code := runCLI(t, "-C", dir, "novel", "revision", "pin", "scene1", "--note", "test pin")
	if code != 0 {
		t.Fatalf("pin exit %d, output=%q", code, out)
	}
	if !strings.Contains(out, "pinned novel scene1") {
		t.Fatalf("output=%q", out)
	}

	out, code = runCLI(t, "-C", dir, "novel", "revision", "list", "scene1")
	if code != 0 || !strings.Contains(out, "current:") {
		t.Fatalf("list exit %d, output=%q", code, out)
	}

	out, code = runCLI(t, "-C", dir, "novel", "show", "scene1")
	if code != 0 || !strings.Contains(out, "pinned_commit:") {
		t.Fatalf("show exit %d, output=%q", code, out)
	}
}