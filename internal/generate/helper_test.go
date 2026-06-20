package generate

import (
	"os"
	"path/filepath"
	"testing"

	"novel-logic/internal/project"
)

func minimalProject(t *testing.T) *project.Data {
	t.Helper()
	dir := t.TempDir()
	files := map[string]string{
		project.FileProject: `title: test
time_order: [t1, t2, t3, t4]
`,
		project.FileThings: `- id: hero
  tags: [character]
  scopes: [plot]
- id: ally
  tags: [character]
  scopes: [plot]
`,
		project.FileTimes: `- id: t1
- id: t2
- id: t3
- id: t4
`,
		project.FileScenes: `- id: scene1
  summary: open
  time_start: t1
  time_end: t2
- id: scene2
  summary: fork
  time_start: t2
  time_end: t4
`,
		project.FileBranches: `- id: main
  label: main line
`,
		project.FileForks:  "[]\n",
		project.FileMerges: "[]\n",
		project.FileFacts: `- id: fact1
  kind: state
  thing: hero
  pred: start
  scope: plot
`,
		project.FileActions: `- id: act1
  thing: hero
  from: start
  to: mid
  at: t2
  scope: plot
`,
		project.FileRules:  "[]\n",
		project.FileNovels: "[]\n",
	}
	for name, content := range files {
		path := filepath.Join(dir, name)
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	d, err := project.Load(dir)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	return d
}