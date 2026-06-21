package cli

import (
	"strings"
	"testing"
)

func TestShowListCommands(t *testing.T) {
	dir := writeCLIProject(t)
	cases := []struct {
		name    string
		args    []string
		contain string
	}{
		{"thing list", []string{"thing", "list"}, "hero"},
		{"fact list", []string{"fact", "list"}, "fact1"},
		{"action list", []string{"action", "list"}, "act1"},
		{"rule list", []string{"rule", "list"}, ""},
		{"time list", []string{"time", "list"}, "t1"},
		{"scene list", []string{"scene", "list"}, "scene1"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			out, code := runCLI(t, append([]string{"-C", dir}, tc.args...)...)
			if code != 0 {
				t.Fatalf("exit %d, output=%q", code, out)
			}
			if tc.contain != "" && !strings.Contains(out, tc.contain) {
				t.Fatalf("output=%q, want substring %q", out, tc.contain)
			}
			if tc.contain == "" && strings.TrimSpace(out) != "" {
				t.Fatalf("output=%q, want empty", out)
			}
		})
	}
}

func TestShowDetailCommands(t *testing.T) {
	dir := writeCLIProject(t)
	cases := []struct {
		name    string
		args    []string
		contain string
	}{
		{"thing show", []string{"thing", "show", "hero"}, "id: hero"},
		{"fact show", []string{"fact", "show", "fact1"}, "kind: state"},
		{"action show", []string{"action", "show", "act1"}, "at: t2"},
		{"scene show", []string{"scene", "show", "scene1"}, "id: scene1"},
		{"plot show", []string{"plot", "show"}, "scene1"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			out, code := runCLI(t, append([]string{"-C", dir}, tc.args...)...)
			if code != 0 {
				t.Fatalf("exit %d, output=%q", code, out)
			}
			if !strings.Contains(out, tc.contain) {
				t.Fatalf("output=%q", out)
			}
		})
	}
}

func TestRuleShow(t *testing.T) {
	dir := writeCLIProject(t)
	if _, code := runCLI(t, "-C", dir, "rule", "add",
		"--kind", "forbid-state", "--thing", "hero", "--pred", "evil"); code != 0 {
		t.Fatal("rule add failed")
	}
	out, code := runCLI(t, "-C", dir, "rule", "show", "rule1")
	if code != 0 || !strings.Contains(out, "kind: forbid-state") {
		t.Fatalf("exit %d, output=%q", code, out)
	}
	out, code = runCLI(t, "-C", dir, "rule", "list")
	if code != 0 || !strings.Contains(out, "rule1") {
		t.Fatalf("list exit %d, output=%q", code, out)
	}
}

func TestFactListFilters(t *testing.T) {
	dir := writeCLIProject(t)
	out, code := runCLI(t, "-C", dir, "fact", "list", "--kind", "state", "--thing", "hero")
	if code != 0 || !strings.Contains(out, "fact1") {
		t.Fatalf("exit %d, output=%q", code, out)
	}
	out, code = runCLI(t, "-C", dir, "fact", "list", "--kind", "fixed")
	if code != 0 || !strings.Contains(out, "(no matches)") {
		t.Fatalf("exit %d, output=%q", code, out)
	}
}

func TestThingListTagFilter(t *testing.T) {
	dir := writeCLIProject(t)
	out, code := runCLI(t, "-C", dir, "thing", "list", "--tag", "character")
	if code != 0 || !strings.Contains(out, "hero") {
		t.Fatalf("exit %d, output=%q", code, out)
	}
	out, code = runCLI(t, "-C", dir, "thing", "list", "--tag", "nope")
	if code != 0 || !strings.Contains(out, "(no matches)") {
		t.Fatalf("exit %d, output=%q", code, out)
	}
}

func TestSceneShowWithNovelLayer(t *testing.T) {
	dir := writeCLIProject(t)
	if _, code := runCLI(t, "-C", dir, "novel", "add", "scene1", "--init"); code != 0 {
		t.Fatal("novel add failed")
	}
	if _, code := runCLI(t, "-C", dir, "fact", "add",
		"--kind", "state", "--thing", "hero", "--pred", "calm", "--scope", "novel:scene1"); code != 0 {
		t.Fatal("fact add failed")
	}
	out, code := runCLI(t, "-C", dir, "scene", "show", "scene1")
	if code != 0 || !strings.Contains(out, "related_things") {
		t.Fatalf("exit %d, output=%q", code, out)
	}
	if !strings.Contains(out, "layers: novel(facts=1") {
		t.Fatalf("output=%q", out)
	}
}

func TestTimelineShowsBranchActions(t *testing.T) {
	dir := writeCLIProject(t)
	out, code := runCLI(t, "-C", dir, "timeline", "--branch", "main", "--verbose")
	if code != 0 {
		t.Fatalf("exit %d, output=%q", code, out)
	}
	if !strings.Contains(out, "act1") || !strings.Contains(out, "facts:") {
		t.Fatalf("output=%q", out)
	}
}

func TestStatusAndTemplateList(t *testing.T) {
	dir := writeCLIProject(t)
	out, code := runCLI(t, "-C", dir, "status")
	if code != 0 || !strings.Contains(out, "entities:") {
		t.Fatalf("status exit %d, output=%q", code, out)
	}
	out, code = runCLI(t, "template", "list")
	if code != 0 {
		t.Fatalf("template list exit %d", code)
	}
	for _, tpl := range []string{"default", "momotaro"} {
		if !strings.Contains(out, tpl) {
			t.Fatalf("output=%q, missing %q", out, tpl)
		}
	}
}