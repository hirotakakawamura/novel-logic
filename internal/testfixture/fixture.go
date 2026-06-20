package testfixture

import (
	"os"
	"path/filepath"
	"testing"

	"novel-logic/internal/project"
)

// WriteFiles writes YAML fixture files into dir.
func WriteFiles(t *testing.T, dir string, files map[string]string) {
	t.Helper()
	for name, content := range files {
		path := filepath.Join(dir, name)
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
}

// MinimalFiles returns a loadable project with two scenes and one action on main.
func MinimalFiles() map[string]string {
	return map[string]string{
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
}

// ValidateFiles returns a smaller project used by validate integration tests.
func ValidateFiles() map[string]string {
	return map[string]string{
		project.FileProject: `title: fixture
time_order: [t1, t2, t3]
`,
		project.FileThings: `- id: hero
  tags: [character]
  scopes: [plot]
`,
		project.FileTimes: `- id: t1
- id: t2
- id: t3
`,
		project.FileScenes: `- id: scene1
  summary: one
  time_start: t1
  time_end: t3
`,
		project.FileBranches: `- id: main
  label: main
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
}

// LoadMinimal writes MinimalFiles to a temp dir and loads the project.
func LoadMinimal(t *testing.T) *project.Data {
	t.Helper()
	dir := t.TempDir()
	WriteFiles(t, dir, MinimalFiles())
	d, err := project.Load(dir)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	return d
}

// WriteMinimalDir writes MinimalFiles to a temp dir and returns its path.
func WriteMinimalDir(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	WriteFiles(t, dir, MinimalFiles())
	return dir
}

// LoadValidate writes ValidateFiles to a temp dir, adds branch_a, and loads.
func LoadValidate(t *testing.T) *project.Data {
	t.Helper()
	dir := t.TempDir()
	WriteFiles(t, dir, ValidateFiles())
	d, err := project.Load(dir)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if err := d.AddBranch("branch_a", "alt", project.MainBranch, "", ""); err != nil {
		t.Fatal(err)
	}
	return d
}