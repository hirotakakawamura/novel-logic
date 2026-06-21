package project

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const NovelBodyExt = ".txt"

// DefaultNovelBodyPath returns novels/<branch>/<scene_id>.txt
func DefaultNovelBodyPath(sceneID, branch string) string {
	branch = NormalizeBranch(branch)
	return filepath.Join(DirNovels, branch, sceneID+NovelBodyExt)
}

// NovelKey is the uniqueness key for novel metadata.
func NovelKey(sceneID, branch string) string {
	return fmt.Sprintf("%s|%s", sceneID, NormalizeBranch(branch))
}

func (d *Data) FindNovel(sceneID, branch string) (*NovelMeta, int) {
	key := NovelKey(sceneID, branch)
	for i := range d.Novels {
		if NovelKey(d.Novels[i].SceneID, d.Novels[i].Branch) == key {
			return &d.Novels[i], i
		}
	}
	return nil, -1
}

func (d *Data) sceneByID(sceneID string) (*Scene, error) {
	for i := range d.Scenes {
		if d.Scenes[i].ID == sceneID {
			return &d.Scenes[i], nil
		}
	}
	return nil, fmt.Errorf("unknown scene %q", sceneID)
}

func validateNovelBodyPath(root, bodyPath string) error {
	if bodyPath == "" {
		return fmt.Errorf("body path is required")
	}
	if filepath.IsAbs(bodyPath) {
		return fmt.Errorf("body path must be relative to project root")
	}
	clean := filepath.Clean(bodyPath)
	if clean == "." || strings.HasPrefix(clean, "..") {
		return fmt.Errorf("body path must stay inside project root")
	}
	return nil
}

func duplicateNovelError(sceneID, branch string) error {
	branch = NormalizeBranch(branch)
	return registrationErrorf(
		"novel for scene %q on branch %q already registered; use: novel-logic novel update %s --branch %s",
		sceneID, branch, sceneID, branch,
	)
}

// AddNovel registers novel metadata pointing at a git-managed body file.
func (d *Data) AddNovel(sceneID, branch, bodyPath string, initFile bool) error {
	if sceneID == "" {
		return fmt.Errorf("scene id is required")
	}
	branch = NormalizeBranch(branch)
	if err := d.validateBranchRef(branch); err != nil {
		return err
	}
	if _, idx := d.FindNovel(sceneID, branch); idx >= 0 {
		return duplicateNovelError(sceneID, branch)
	}
	scene, err := d.sceneByID(sceneID)
	if err != nil {
		return err
	}
	if bodyPath == "" {
		bodyPath = DefaultNovelBodyPath(sceneID, branch)
	}
	if err := validateNovelBodyPath(d.Root, bodyPath); err != nil {
		return err
	}
	abs := filepath.Join(d.Root, bodyPath)
	if initFile {
		if err := os.MkdirAll(filepath.Dir(abs), 0o755); err != nil {
			return err
		}
		if _, err := os.Stat(abs); os.IsNotExist(err) {
			if err := os.WriteFile(abs, nil, 0o644); err != nil {
				return err
			}
		} else if err != nil {
			return err
		}
	} else if _, err := os.Stat(abs); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("body file %q not found; create it in git or use --init", bodyPath)
		}
		return err
	}
	d.Novels = append(d.Novels, NovelMeta{
		SceneID:   sceneID,
		Branch:    branch,
		TimeStart: scene.TimeStart,
		TimeEnd:   scene.TimeEnd,
		BodyPath:  bodyPath,
	})
	return nil
}

// UpdateNovel updates novel metadata. Body text is never written by the tool.
func (d *Data) UpdateNovel(sceneID, branch, bodyPath string, bodyPathChanged bool) error {
	if sceneID == "" {
		return fmt.Errorf("scene id is required")
	}
	branch = NormalizeBranch(branch)
	n, idx := d.FindNovel(sceneID, branch)
	if n == nil {
		return fmt.Errorf("novel for scene %q on branch %q not found", sceneID, branch)
	}
	scene, err := d.sceneByID(sceneID)
	if err != nil {
		return err
	}
	nextPath := n.BodyPath
	if bodyPathChanged {
		if bodyPath == "" {
			return fmt.Errorf("body path is required")
		}
		if err := validateNovelBodyPath(d.Root, bodyPath); err != nil {
			return err
		}
		abs := filepath.Join(d.Root, bodyPath)
		if _, err := os.Stat(abs); err != nil {
			if os.IsNotExist(err) {
				return fmt.Errorf("body file %q not found", bodyPath)
			}
			return err
		}
		nextPath = bodyPath
	}
	updated := *n
	updated.TimeStart = scene.TimeStart
	updated.TimeEnd = scene.TimeEnd
	updated.BodyPath = nextPath
	d.Novels[idx] = updated
	return nil
}

// RemoveNovel drops novel metadata for a scene on a branch.
func (d *Data) RemoveNovel(sceneID, branch string, keepBody bool) error {
	if sceneID == "" {
		return fmt.Errorf("scene id is required")
	}
	branch = NormalizeBranch(branch)
	var bodyPath string
	found := false
	out := d.Novels[:0]
	for _, n := range d.Novels {
		if n.SceneID == sceneID && NormalizeBranch(n.Branch) == branch {
			found = true
			bodyPath = n.BodyPath
			continue
		}
		out = append(out, n)
	}
	if !found {
		return fmt.Errorf("novel for scene %q on branch %q not found", sceneID, branch)
	}
	d.Novels = out
	if !keepBody && bodyPath != "" {
		p := filepath.Join(d.Root, bodyPath)
		if err := os.Remove(p); err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("remove body file: %w", err)
		}
	}
	return nil
}

// NovelBodyIssues reports missing body files for registered novels.
func NovelBodyIssues(d *Data) []string {
	var issues []string
	for _, n := range d.Novels {
		if n.BodyPath == "" {
			issues = append(issues, fmt.Sprintf("novel %q branch %q: empty body_path", n.SceneID, NormalizeBranch(n.Branch)))
			continue
		}
		expected := DefaultNovelBodyPath(n.SceneID, n.Branch)
		if filepath.ToSlash(n.BodyPath) != filepath.ToSlash(expected) {
			issues = append(issues, fmt.Sprintf("novel %q branch %q: body path %q (expected %q)", n.SceneID, NormalizeBranch(n.Branch), n.BodyPath, expected))
		}
		p := filepath.Join(d.Root, n.BodyPath)
		if _, err := os.Stat(p); err != nil {
			if os.IsNotExist(err) {
				issues = append(issues, fmt.Sprintf("novel %q branch %q: body file missing at %s", n.SceneID, NormalizeBranch(n.Branch), n.BodyPath))
			} else {
				issues = append(issues, fmt.Sprintf("novel %q branch %q: cannot read body file %s: %v", n.SceneID, NormalizeBranch(n.Branch), n.BodyPath, err))
			}
		}
	}
	return issues
}