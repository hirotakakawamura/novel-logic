package project

import "time"

// ProjectMeta lives in project.yaml.
type ProjectMeta struct {
	Title     string        `yaml:"title"`
	CreatedAt time.Time     `yaml:"created_at,omitempty"`
	TimeOrder []string      `yaml:"time_order"`
	LastCheck *CheckResult  `yaml:"last_check,omitempty"`
}

type CheckResult struct {
	At      time.Time `yaml:"at"`
	Success bool      `yaml:"success"`
	Stage1  bool      `yaml:"stage1"`
	Stage2  bool      `yaml:"stage2"`
	Message string    `yaml:"message,omitempty"`
}

// Plot is stored in plot.yaml.
type Plot struct {
	Summary string `yaml:"summary"`
}

// Thing is stored in things.yaml.
type Thing struct {
	ID     string   `yaml:"id"`
	Name   string   `yaml:"name,omitempty"`
	Tags   []string `yaml:"tags"`
	Scopes []string `yaml:"scopes"`
}

// Scene is stored in scenes.yaml.
type Scene struct {
	ID        string `yaml:"id"`
	Summary   string `yaml:"summary"`
	TimeStart string `yaml:"time_start"`
	TimeEnd   string `yaml:"time_end"`
}

// TimeEntry is stored in times.yaml (existence registry).
type TimeEntry struct {
	ID string `yaml:"id"`
}

// FactKind is fixed_fact or state.
type FactKind string

const (
	FactFixed FactKind = "fixed"
	FactState FactKind = "state"
)

// Branch is stored in branches.yaml.
type Branch struct {
	ID        string `yaml:"id"`
	Label     string `yaml:"label,omitempty"`
	Parent    string `yaml:"parent,omitempty"`
	ViaFork   string `yaml:"via_fork,omitempty"`
	ViaAction string `yaml:"via_action,omitempty"`
}

// ForkChoice maps a fork action to a child branch.
type ForkChoice struct {
	Action string `yaml:"action"`
	Branch string `yaml:"branch"`
}

// Fork is stored in forks.yaml.
type Fork struct {
	ID           string       `yaml:"id"`
	ParentBranch string       `yaml:"parent_branch"`
	At           string       `yaml:"at"`
	Scope        string       `yaml:"scope"`
	Choices      []ForkChoice `yaml:"choices"`
}

// MergeChoice maps a from-branch to its merge action.
type MergeChoice struct {
	Branch string `yaml:"branch"`
	Action string `yaml:"action"`
}

// Merge is stored in merges.yaml.
type Merge struct {
	ID         string        `yaml:"id"`
	At         string        `yaml:"at"`
	Scope      string        `yaml:"scope"`
	IntoBranch string        `yaml:"into_branch"`
	Choices    []MergeChoice `yaml:"choices"`
}

// Fact is stored in facts.yaml.
type Fact struct {
	ID     string   `yaml:"id"`
	Kind   FactKind `yaml:"kind"`
	Thing  string   `yaml:"thing"`
	Pred   string   `yaml:"pred"`
	Scope  string   `yaml:"scope"`
	Branch string   `yaml:"branch,omitempty"`
}

// Action is stored in actions.yaml.
type Action struct {
	ID     string `yaml:"id"`
	Thing  string `yaml:"thing"`
	From   string `yaml:"from,omitempty"`
	To     string `yaml:"to"`
	At     string `yaml:"at"`
	Scope  string `yaml:"scope"`
	Label  string `yaml:"label,omitempty"`
	Branch string `yaml:"branch,omitempty"`
}

// RuleKind matches CLI rule kinds.
type RuleKind string

const (
	RuleForbidState       RuleKind = "forbid-state"
	RuleForbidTransition  RuleKind = "forbid-transition"
)

// Rule is stored in rules.yaml.
type Rule struct {
	ID     string   `yaml:"id"`
	Kind   RuleKind `yaml:"kind"`
	Thing  string   `yaml:"thing,omitempty"`
	Pred   string   `yaml:"pred,omitempty"`
	From   string   `yaml:"from,omitempty"`
	To     string   `yaml:"to,omitempty"`
	Branch string   `yaml:"branch,omitempty"`
}

// NovelRevision records a git commit pinned to a novel body file.
type NovelRevision struct {
	Revision   string    `yaml:"revision"`
	Short      string    `yaml:"short,omitempty"`
	Branch     string    `yaml:"branch,omitempty"`
	RecordedAt time.Time `yaml:"recorded_at,omitempty"`
	Note       string    `yaml:"note,omitempty"`
	Dirty      bool      `yaml:"dirty,omitempty"`
}

// NovelMeta is stored in novels.yaml (keyed by scene_id + branch).
type NovelMeta struct {
	SceneID   string          `yaml:"scene_id"`
	Branch    string          `yaml:"branch,omitempty"`
	TimeStart string          `yaml:"time_start"`
	TimeEnd   string          `yaml:"time_end"`
	BodyPath  string          `yaml:"body_path"`
	Revision  string          `yaml:"revision,omitempty"`
	Revisions []NovelRevision `yaml:"revisions,omitempty"`
}

// Data is the in-memory representation of a work project.
type Data struct {
	Root      string
	Meta      ProjectMeta
	Plot      Plot
	Things    []Thing
	Scenes    []Scene
	Times     []TimeEntry
	Branches  []Branch
	Forks     []Fork
	Merges    []Merge
	Facts     []Fact
	Actions   []Action
	Rules     []Rule
	Novels    []NovelMeta
}