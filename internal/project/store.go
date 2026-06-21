package project

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const (
	FileProject  = "project.yaml"
	FilePlot     = "plot.yaml"
	FileThings   = "things.yaml"
	FileScenes   = "scenes.yaml"
	FileTimes    = "times.yaml"
	FileBranches = "branches.yaml"
	FileForks    = "forks.yaml"
	FileMerges   = "merges.yaml"
	FileFacts    = "facts.yaml"
	FileActions  = "actions.yaml"
	FileRules    = "rules.yaml"
	FileNovels   = "novels.yaml"
	DirNovels    = "novels"
	DirLogic     = "logic"
)

func Load(root string) (*Data, error) {
	root, err := filepath.Abs(root)
	if err != nil {
		return nil, err
	}
	if _, err := os.Stat(filepath.Join(root, FileProject)); err != nil {
		return nil, fmt.Errorf("not a novel-logic project (missing %s): %w", FileProject, err)
	}
	// Optional YAML files are read best-effort; missing paths yield empty slices.
	// ensureMainBranch adds main when branches.yaml is absent. See doctor recommended_missing.
	d := &Data{Root: root}
	if err := readYAML(filepath.Join(root, FileProject), &d.Meta); err != nil {
		return nil, fmt.Errorf("read %s: %w", FileProject, err)
	}
	_ = readYAML(filepath.Join(root, FilePlot), &d.Plot)
	_ = readYAML(filepath.Join(root, FileThings), &d.Things)
	_ = readYAML(filepath.Join(root, FileScenes), &d.Scenes)
	_ = readYAML(filepath.Join(root, FileTimes), &d.Times)
	_ = readYAML(filepath.Join(root, FileBranches), &d.Branches)
	_ = readYAML(filepath.Join(root, FileForks), &d.Forks)
	_ = readYAML(filepath.Join(root, FileMerges), &d.Merges)
	_ = readYAML(filepath.Join(root, FileFacts), &d.Facts)
	_ = readYAML(filepath.Join(root, FileActions), &d.Actions)
	_ = readYAML(filepath.Join(root, FileRules), &d.Rules)
	_ = readYAML(filepath.Join(root, FileNovels), &d.Novels)
	normalizeLoadedData(d)
	return d, nil
}

func normalizeLoadedData(d *Data) {
	for i := range d.Facts {
		d.Facts[i].Branch = NormalizeBranch(d.Facts[i].Branch)
	}
	for i := range d.Actions {
		d.Actions[i].Branch = NormalizeBranch(d.Actions[i].Branch)
	}
	for i := range d.Rules {
		d.Rules[i].Branch = NormalizeBranch(d.Rules[i].Branch)
	}
	for i := range d.Novels {
		d.Novels[i].Branch = NormalizeBranch(d.Novels[i].Branch)
		if d.Novels[i].BodyPath == "" {
			d.Novels[i].BodyPath = DefaultNovelBodyPath(d.Novels[i].SceneID, d.Novels[i].Branch)
		} else {
			d.Novels[i].BodyPath = migrateNovelBodyPath(d.Novels[i].BodyPath, d.Novels[i].SceneID, d.Novels[i].Branch)
		}
	}
	for i := range d.Forks {
		d.Forks[i].ParentBranch = NormalizeBranch(d.Forks[i].ParentBranch)
	}
	for i := range d.Merges {
		d.Merges[i].IntoBranch = NormalizeBranch(d.Merges[i].IntoBranch)
	}
	ensureMainBranch(d)
}

func ensureMainBranch(d *Data) {
	if _, idx := d.FindBranchDef(MainBranch); idx >= 0 {
		return
	}
	d.Branches = append([]Branch{{ID: MainBranch, Label: "本線"}}, d.Branches...)
}

// migrateNovelBodyPath rewrites legacy novels/<scene>.txt to novels/<branch>/<scene>.txt.
func migrateNovelBodyPath(path, sceneID, branch string) string {
	branch = NormalizeBranch(branch)
	legacy := filepath.Join(DirNovels, sceneID+NovelBodyExt)
	if path == legacy || path == filepath.ToSlash(legacy) {
		return DefaultNovelBodyPath(sceneID, branch)
	}
	return path
}

func Save(d *Data) error {
	if err := writeYAML(filepath.Join(d.Root, FileProject), d.Meta); err != nil {
		return err
	}
	if err := writeYAML(filepath.Join(d.Root, FilePlot), d.Plot); err != nil {
		return err
	}
	if err := writeYAML(filepath.Join(d.Root, FileThings), d.Things); err != nil {
		return err
	}
	if err := writeYAML(filepath.Join(d.Root, FileScenes), d.Scenes); err != nil {
		return err
	}
	if err := writeYAML(filepath.Join(d.Root, FileTimes), d.Times); err != nil {
		return err
	}
	if err := writeYAML(filepath.Join(d.Root, FileBranches), d.Branches); err != nil {
		return err
	}
	if err := writeYAML(filepath.Join(d.Root, FileForks), d.Forks); err != nil {
		return err
	}
	if err := writeYAML(filepath.Join(d.Root, FileMerges), d.Merges); err != nil {
		return err
	}
	if err := writeYAML(filepath.Join(d.Root, FileFacts), d.Facts); err != nil {
		return err
	}
	if err := writeYAML(filepath.Join(d.Root, FileActions), d.Actions); err != nil {
		return err
	}
	if err := writeYAML(filepath.Join(d.Root, FileRules), d.Rules); err != nil {
		return err
	}
	return writeYAML(filepath.Join(d.Root, FileNovels), d.Novels)
}

func readYAML(path string, out any) error {
	b, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	if len(b) == 0 {
		return nil
	}
	return yaml.Unmarshal(b, out)
}

func writeYAML(path string, in any) error {
	b, err := yaml.Marshal(in)
	if err != nil {
		return err
	}
	return os.WriteFile(path, b, 0o644)
}

func (d *Data) ThingIDs() map[string]bool {
	m := make(map[string]bool, len(d.Things))
	for _, t := range d.Things {
		m[t.ID] = true
	}
	return m
}

func (d *Data) Preds() map[string]bool {
	m := make(map[string]bool)
	for _, f := range d.Facts {
		m[f.Pred] = true
	}
	for _, a := range d.Actions {
		if a.From != "" {
			m[a.From] = true
		}
		m[a.To] = true
	}
	for _, r := range d.Rules {
		if r.Pred != "" {
			m[r.Pred] = true
		}
		if r.From != "" {
			m[r.From] = true
		}
		if r.To != "" {
			m[r.To] = true
		}
	}
	return m
}

func (d *Data) FindThing(id string) (*Thing, int) {
	for i := range d.Things {
		if d.Things[i].ID == id {
			return &d.Things[i], i
		}
	}
	return nil, -1
}

func (d *Data) TimeIndex(id string) int {
	for i, t := range d.Meta.TimeOrder {
		if t == id {
			return i
		}
	}
	return -1
}

func (d *Data) TimeLE(a, b string) bool {
	ia, ib := d.TimeIndex(a), d.TimeIndex(b)
	if ia < 0 || ib < 0 {
		return false
	}
	return ia <= ib
}