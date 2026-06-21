package cli

import (
	"bytes"
	"errors"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"novel-logic/internal/project"
	"novel-logic/internal/testfixture"
)

func resetCLIGlobals() {
	projectPath = "."
	quiet = false
	verbose = false
	checkQuick = false
	checkNoGenerate = false
	checkJobs = 0
	// rootCmd is a package-global cobra tree; flags persist across Execute calls unless reset.
	resetCLIFlags(rootCmd)
}

func resetCLIFlags(cmd *cobra.Command) {
	cmd.Flags().VisitAll(resetFlagValue)
	cmd.PersistentFlags().VisitAll(resetFlagValue)
	for _, c := range cmd.Commands() {
		resetCLIFlags(c)
	}
}

func resetFlagValue(f *pflag.Flag) {
	f.Changed = false
	if sv, ok := f.Value.(pflag.SliceValue); ok {
		_ = sv.Replace(nil)
		return
	}
	_ = f.Value.Set(f.DefValue)
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
	return testfixture.WriteMinimalDir(t)
}

func gitInit(t *testing.T, dir string) {
	t.Helper()
	runGitCLI(t, dir, "init")
	runGitCLI(t, dir, "config", "user.email", "test@example.com")
	runGitCLI(t, dir, "config", "user.name", "Test User")
	gitCommitAll(t, dir, "init")
}

func gitCommitAll(t *testing.T, dir, msg string) {
	t.Helper()
	runGitCLI(t, dir, "add", ".")
	runGitCLI(t, dir, "commit", "-m", msg)
}

func runGitCLI(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %s: %v\n%s", strings.Join(args, " "), err, out)
	}
}

func lastActionID(t *testing.T, dir string) string {
	t.Helper()
	d, err := project.Load(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(d.Actions) == 0 {
		t.Fatal("no actions")
	}
	return d.Actions[len(d.Actions)-1].ID
}

func mustOK(t *testing.T, dir string, args ...string) {
	t.Helper()
	if _, code := runCLI(t, append([]string{"-C", dir}, args...)...); code != 0 {
		t.Fatalf("command failed: %v (exit %d)", args, code)
	}
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
