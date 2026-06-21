package generate

import (
	"strings"
	"testing"

	"novel-logic/internal/project"
)

func TestScopeExpr(t *testing.T) {
	cases := []struct {
		scope string
		want  string
	}{
		{"", "Scope.plot"},
		{"plot", "Scope.plot"},
		{"novel:scene1", "Scope.novel_scene1"},
		{"bad", "Scope.plot"},
	}
	for _, tc := range cases {
		if got := scopeExpr(tc.scope); got != tc.want {
			t.Errorf("scopeExpr(%q) = %q, want %q", tc.scope, got, tc.want)
		}
	}
}

func TestLeanIdentEdgeCases(t *testing.T) {
	cases := []struct {
		in, want string
	}{
		{"", "none_"},
		{"9start", "t9start"},
		{"a-b", "a_b"},
		{"scene.1", "scene_1"},
	}
	for _, tc := range cases {
		if got := leanIdent(tc.in); got != tc.want {
			t.Errorf("leanIdent(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

func TestIsASCIIIdent(t *testing.T) {
	if !isASCIIIdent("valid_id2") {
		t.Fatal("expected valid")
	}
	if isASCIIIdent("") || isASCIIIdent("桃太郎") || isASCIIIdent("bad-id") {
		t.Fatal("expected invalid")
	}
}

func TestLeanPredRomanizationAndHash(t *testing.T) {
	if got := leanPred("人間"); got != "ningen" {
		t.Fatalf("leanPred(人間) = %q", got)
	}
	hash := leanPred("未知の述語")
	if hash == "" || hash == "未知の述語" {
		t.Fatalf("leanPred hash = %q", hash)
	}
	if got := leanPred(""); got != "none_" {
		t.Fatalf("leanPred empty = %q", got)
	}
}

func TestActionDeclExprNovelScope(t *testing.T) {
	d := minimalProject(t)
	act, err := d.AddAction("hero", "start", "novel_end", "t2", "novel:scene1", "", project.MainBranch)
	if err != nil {
		t.Fatal(err)
	}
	expr := actionDeclExpr(act)
	if !strings.Contains(expr, "Scope.novel_scene1") || !strings.Contains(expr, "novel_end") {
		t.Fatalf("expr = %q", expr)
	}
}