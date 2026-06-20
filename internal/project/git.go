package project

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
)

// GitFileState describes the git revision of a tracked file.
type GitFileState struct {
	Revision string
	Short    string
	Branch   string
	Dirty    bool
}

func gitRun(dir string, args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	out, err := cmd.Output()
	if err != nil {
		if ee, ok := err.(*exec.ExitError); ok {
			msg := strings.TrimSpace(string(ee.Stderr))
			if msg != "" {
				return "", fmt.Errorf("%s", msg)
			}
		}
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

func isGitRepo(root string) bool {
	_, err := gitRun(root, "rev-parse", "--git-dir")
	return err == nil
}

// IsGitRepo reports whether root is inside a git work tree.
func IsGitRepo(root string) bool {
	return isGitRepo(root)
}

func fileTrackedInGit(root, relPath string) bool {
	_, err := gitRun(root, "ls-files", "--error-unmatch", "--", relPath)
	return err == nil
}

// ResolveGitFileState returns the latest commit touching relPath and whether the
// working tree copy differs from that commit.
func ResolveGitFileState(root, relPath string) (GitFileState, error) {
	if !isGitRepo(root) {
		return GitFileState{}, fmt.Errorf("project is not a git repository")
	}
	relPath = filepath.ToSlash(relPath)
	if !fileTrackedInGit(root, relPath) {
		return GitFileState{}, fmt.Errorf("file %q is not tracked by git (commit it first)", relPath)
	}
	rev, err := gitRun(root, "log", "-1", "--format=%H", "--", relPath)
	if err != nil || rev == "" {
		return GitFileState{}, fmt.Errorf("no git history for %q", relPath)
	}
	short, err := gitRun(root, "rev-parse", "--short", rev)
	if err != nil {
		short = rev
		if len(short) > 12 {
			short = short[:12]
		}
	}
	branch, _ := gitRun(root, "rev-parse", "--abbrev-ref", "HEAD")
	status, _ := gitRun(root, "status", "--porcelain", "--", relPath)
	dirty := status != ""
	return GitFileState{
		Revision: rev,
		Short:    short,
		Branch:   branch,
		Dirty:    dirty,
	}, nil
}

// ResolveGitRevisionAt returns git metadata for an explicit commit and path.
func ResolveGitRevisionAt(root, relPath, revision string) (GitFileState, error) {
	if !isGitRepo(root) {
		return GitFileState{}, fmt.Errorf("project is not a git repository")
	}
	relPath = filepath.ToSlash(relPath)
	if _, err := gitRun(root, "cat-file", "-e", revision+":"+relPath); err != nil {
		return GitFileState{}, fmt.Errorf("file %q not found at revision %s", relPath, revision)
	}
	short, err := gitRun(root, "rev-parse", "--short", revision)
	if err != nil {
		short = revision
	}
	branch, _ := gitRun(root, "name-rev", "--name-only", "--no-undefined", revision)
	branch = strings.TrimSuffix(branch, "~0")
	branch = strings.TrimSuffix(branch, "^0")
	return GitFileState{
		Revision: revision,
		Short:    short,
		Branch:   branch,
	}, nil
}

func currentHEAD(root string) (string, error) {
	return gitRun(root, "rev-parse", "HEAD")
}

func workingTreeDirty(root, relPath string) bool {
	status, err := gitRun(root, "status", "--porcelain", "--", filepath.ToSlash(relPath))
	return err == nil && status != ""
}

func revisionsEqual(a, b string) bool {
	return a == b || strings.HasPrefix(a, b) || strings.HasPrefix(b, a)
}