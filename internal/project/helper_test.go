package project

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTestFile(t *testing.T, dir, name, content string) {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

// newTestProject creates a minimal loadable project in a temp directory.
func newTestProject(t *testing.T) *Data {
	t.Helper()
	dir := t.TempDir()
	writeTestFile(t, dir, FileProject, `title: test
time_order: [t1, t2, t3, t4]
`)
	writeTestFile(t, dir, FileThings, `- id: hero
  tags: [character]
  scopes: [plot]
- id: ally
  tags: [character]
  scopes: [plot]
`)
	writeTestFile(t, dir, FileTimes, `- id: t1
- id: t2
- id: t3
- id: t4
`)
	writeTestFile(t, dir, FileScenes, `- id: scene1
  summary: open
  time_start: t1
  time_end: t2
- id: scene2
  summary: fork
  time_start: t2
  time_end: t4
`)
	writeTestFile(t, dir, FileBranches, `- id: main
  label: main line
`)
	writeTestFile(t, dir, FileForks, "[]\n")
	writeTestFile(t, dir, FileMerges, "[]\n")
	writeTestFile(t, dir, FileFacts, `- id: fact1
  kind: state
  thing: hero
  pred: start
  scope: plot
`)
	writeTestFile(t, dir, FileActions, `- id: act1
  thing: hero
  from: start
  to: mid
  at: t2
  scope: plot
`)
	writeTestFile(t, dir, FileRules, "[]\n")
	writeTestFile(t, dir, FileNovels, "[]\n")

	d, err := Load(dir)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	return d
}

func actionIDs(actions []Action) []string {
	ids := make([]string, len(actions))
	for i, a := range actions {
		ids[i] = a.ID
	}
	return ids
}

func containsID(ids []string, want string) bool {
	for _, id := range ids {
		if id == want {
			return true
		}
	}
	return false
}