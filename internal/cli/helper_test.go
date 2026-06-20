package cli

import (
	"bytes"
	"errors"
	"io"
	"os"
	"path/filepath"
	"testing"

	"novel-logic/internal/project"
)

func resetCLIGlobals() {
	projectPath = "."
	quiet = false
	verbose = false
	checkQuick = false
	checkNoGenerate = false
	checkJobs = 0
}

func runCLI(t *testing.T, args ...string) (stdout string, exitCode int) {
	t.Helper()
	resetCLIGlobals()

	stdoutR, stdoutW, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	stderrR, stderrW, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}

	oldStdout := os.Stdout
	oldStderr := os.Stderr
	os.Stdout = stdoutW
	os.Stderr = stderrW

	rootCmd.SetIn(os.Stdin)
	rootCmd.SetArgs(args)
	rootCmd.SilenceErrors = true
	rootCmd.SilenceUsage = true

	execErr := rootCmd.Execute()

	stdoutW.Close()
	stderrW.Close()
	os.Stdout = oldStdout
	os.Stderr = oldStderr

	var combined bytes.Buffer
	_, _ = io.Copy(&combined, stdoutR)
	_, _ = io.Copy(&combined, stderrR)
	_ = stdoutR.Close()
	_ = stderrR.Close()

	exitCode = 0
	if execErr != nil {
		var ee *ExitError
		if errors.As(execErr, &ee) {
			exitCode = ee.Code
		} else {
			exitCode = 1
		}
	}
	return combined.String(), exitCode
}

func writeCLIProject(t *testing.T) string {
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
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return dir
}

func copyWalkthroughProject(t *testing.T) string {
	t.Helper()
	src := filepath.Join("..", "..", "examples", "momotaro-walkthrough")
	dst := t.TempDir()
	if err := copyDir(src, dst); err != nil {
		t.Fatal(err)
	}
	return dst
}

func copyDir(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		if rel == "." {
			return nil
		}
		target := filepath.Join(dst, rel)
		if info.IsDir() {
			return os.MkdirAll(target, info.Mode())
		}
		in, err := os.Open(path)
		if err != nil {
			return err
		}
		defer in.Close()
		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
			return err
		}
		out, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, info.Mode())
		if err != nil {
			return err
		}
		defer out.Close()
		_, err = io.Copy(out, in)
		return err
	})
}