package generate

import (
	"crypto/sha256"
	"embed"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"unicode"

	"novel-logic/internal/project"
)

//go:embed templates/Core.lean
var coreFS embed.FS

const CoreVersion = "0.1.0"

func Run(d *project.Data) error {
	logicDir := filepath.Join(d.Root, project.DirLogic)
	if err := os.MkdirAll(logicDir, 0o755); err != nil {
		return err
	}
	core, err := coreFS.ReadFile("templates/Core.lean")
	if err != nil {
		return fmt.Errorf("read Core.lean: %w", err)
	}
	files := map[string]string{
		"Core.lean":     string(core),
		"Project.lean":  genProject(d),
		"Facts.lean":    genFacts(d),
		"Rules.lean":    genRules(d),
		"Timeline.lean": genTimeline(d),
		"Theorems.lean": genTheorems(d),

		"lakefile.toml": genLakefile(d),
	}
	for name, content := range files {
		if err := os.WriteFile(filepath.Join(logicDir, name), []byte(content), 0o644); err != nil {
			return fmt.Errorf("write %s: %w", name, err)
		}
	}
	return nil
}

func genLakefile(d *project.Data) string {
	return `name = "novel-logic-work"
version = "0.1.0"
defaultTargets = ["Theorems"]

[[lean_lib]]
name = "NovelLogic"
roots = ["Core", "Project", "Facts", "Rules", "Timeline", "Theorems"]
`
}

func genProject(d *project.Data) string {
	ns := namespace(d)
	var b strings.Builder
	fmt.Fprintf(&b, "import Core\n\nnamespace %s\n\n", ns)
	b.WriteString(genInductive("ThingId", thingIDs(d)))
	b.WriteString(genInductive("TimeId", d.Meta.TimeOrder))
	b.WriteString(genInductive("BranchId", branchIDs(d)))
	b.WriteString(genInductive("SceneId", sceneIDs(d)))
	b.WriteString(genInductive("PredId", predIDs(d)))
	b.WriteString("inductive Scope\n  | plot\n")
	for _, sid := range sceneIDs(d) {
		fmt.Fprintf(&b, "  | novel_%s\n", leanIdent(sid))
	}
	fmt.Fprintf(&b, "  deriving DecidableEq, Repr\n\n")
	fmt.Fprintf(&b, "abbrev PlotScope : Scope := Scope.plot\n\n")
	fmt.Fprintf(&b, "def timeOrder : List TimeId := [%s]\n\n", joinIDs(d.Meta.TimeOrder, "TimeId"))
	fmt.Fprintf(&b, "def scopeToScene : Scope → Option SceneId\n")
	fmt.Fprintf(&b, "  | Scope.plot => none\n")
	for _, sid := range sceneIDs(d) {
		ident := leanIdent(sid)
		fmt.Fprintf(&b, "  | Scope.novel_%s => some SceneId.%s\n", ident, ident)
	}
	fmt.Fprintf(&b, "\nend %s\n", ns)
	return b.String()
}

func genFacts(d *project.Data) string {
	ns := namespace(d)
	var b strings.Builder
	fmt.Fprintf(&b, "import Core\nimport Project\n\nnamespace %s\n\nopen NovelLogic\n\n", ns)
	bids := branchIDs(d)
	for _, bid := range bids {
		writeBranchFactLists(&b, d, bid)
	}
	mainIdent := leanIdent(project.MainBranch)
	fmt.Fprintf(&b, "abbrev allFixedFacts := fixedFacts_%s\n\n", mainIdent)
	fmt.Fprintf(&b, "abbrev allStateDecls := stateDecls_%s\n\n", mainIdent)
	for _, bid := range bids {
		ident := leanIdent(bid)
		acts := d.ActiveActions(bid)
		fmt.Fprintf(&b, "def activeActions_%s : List (ActionDecl ThingId PredId TimeId Scope) := [\n", ident)
		for _, a := range acts {
			fmt.Fprintf(&b, "  %s,\n", actionDeclExpr(a))
		}
		fmt.Fprintf(&b, "]\n\n")
		fmt.Fprintf(&b, "def evolveBranch_%s (t : TimeId) (thing : ThingId) : List PredId :=\n", ident)
		fmt.Fprintf(&b, "  predsAt fixedFacts_%s stateDecls_%s activeActions_%s timeOrder t thing\n\n", ident, ident, ident)
	}
	fmt.Fprintf(&b, "end %s\n", ns)
	return b.String()
}

func writeBranchFactLists(b *strings.Builder, d *project.Data, branchID string) {
	ident := leanIdent(branchID)
	facts := d.EffectiveFactsOnBranch(branchID)
	fmt.Fprintf(b, "def fixedFacts_%s : List (FixedFact ThingId PredId Scope) := [\n", ident)
	for _, f := range facts {
		if f.Kind != project.FactFixed {
			continue
		}
		fmt.Fprintf(b, "  ⟨ThingId.%s, PredId.%s, %s⟩,\n", leanIdent(f.Thing), leanPred(f.Pred), scopeExpr(f.Scope))
	}
	fmt.Fprintf(b, "]\n\n")
	fmt.Fprintf(b, "def stateDecls_%s : List (StateDecl ThingId PredId Scope) := [\n", ident)
	for _, f := range facts {
		if f.Kind != project.FactState {
			continue
		}
		fmt.Fprintf(b, "  ⟨ThingId.%s, PredId.%s, %s⟩,\n", leanIdent(f.Thing), leanPred(f.Pred), scopeExpr(f.Scope))
	}
	fmt.Fprintf(b, "]\n\n")
}

func genRules(d *project.Data) string {
	ns := namespace(d)
	var b strings.Builder
	fmt.Fprintf(&b, "import Core\nimport Project\n\nnamespace %s\n\nopen NovelLogic\n\n", ns)
	for _, bid := range branchIDs(d) {
		ident := leanIdent(bid)
		rules := d.EffectiveRulesOnBranch(bid)
		fmt.Fprintf(&b, "def projectRules_%s : Rules ThingId PredId := {\n", ident)
		fmt.Fprintf(&b, "  forbiddenStates := [\n")
		for _, r := range rules {
			if r.Kind == project.RuleForbidState {
				fmt.Fprintf(&b, "    (ThingId.%s, PredId.%s),\n", leanIdent(r.Thing), leanPred(r.Pred))
			}
		}
		fmt.Fprintf(&b, "  ],\n  forbiddenTransitions := [\n")
		for _, r := range rules {
			if r.Kind == project.RuleForbidTransition {
				fmt.Fprintf(&b, "    (PredId.%s, PredId.%s),\n", leanPred(r.From), leanPred(r.To))
			}
		}
		fmt.Fprintf(&b, "  ]\n}\n\n")
	}
	mainIdent := leanIdent(project.MainBranch)
	fmt.Fprintf(&b, "abbrev projectRules : Rules ThingId PredId := projectRules_%s\n\n", mainIdent)
	fmt.Fprintf(&b, "end %s\n", ns)
	return b.String()
}

func genTimeline(d *project.Data) string {
	ns := namespace(d)
	var b strings.Builder
	fmt.Fprintf(&b, "import Core\nimport Project\n\nnamespace %s\n\nopen NovelLogic\n\n", ns)
	fmt.Fprintf(&b, "def sceneWindows : List (SceneWindow SceneId TimeId) := [\n")
	for _, s := range d.Scenes {
		fmt.Fprintf(&b, "  ⟨SceneId.%s, TimeId.%s, TimeId.%s⟩,\n",
			leanIdent(s.ID), leanIdent(s.TimeStart), leanIdent(s.TimeEnd))
	}
	fmt.Fprintf(&b, "]\n\nend %s\n", ns)
	return b.String()
}

func genTheorems(d *project.Data) string {
	ns := namespace(d)
	lastTime := ""
	if n := len(d.Meta.TimeOrder); n > 0 {
		lastTime = d.Meta.TimeOrder[n-1]
	}
	var b strings.Builder
	fmt.Fprintf(&b, "import Core\nimport Project\nimport Facts\nimport Rules\nimport Timeline\n\nnamespace %s\n\n", ns)
	fmt.Fprintf(&b, "open NovelLogic\n\n")
	mainIdent := leanIdent(project.MainBranch)
	fmt.Fprintf(&b, "theorem actions_in_scene_window :\n")
	fmt.Fprintf(&b, "    allActionsInSceneWindows sceneWindows timeOrder activeActions_%s scopeToScene := by\n", mainIdent)
	fmt.Fprintf(&b, "  native_decide\n\n")
	fmt.Fprintf(&b, "theorem no_forbidden_states :\n")
	fmt.Fprintf(&b, "    noForbiddenStatesRegistered projectRules stateDecls_%s := by\n", mainIdent)
	fmt.Fprintf(&b, "  native_decide\n\n")
	fmt.Fprintf(&b, "theorem no_forbidden_transitions :\n")
	fmt.Fprintf(&b, "    allActionsRespectRules projectRules activeActions_%s := by\n", mainIdent)
	fmt.Fprintf(&b, "  native_decide\n\n")
	fmt.Fprintf(&b, "theorem fixed_facts_stable :\n")
	fmt.Fprintf(&b, "    fixedFactsStable fixedFacts_%s stateDecls_%s activeActions_%s timeOrder := by\n", mainIdent, mainIdent, mainIdent)
	fmt.Fprintf(&b, "  native_decide\n\n")
	for _, bid := range branchIDs(d) {
		if project.NormalizeBranch(bid) == project.MainBranch {
			continue
		}
		ident := leanIdent(bid)
		fmt.Fprintf(&b, "theorem actions_in_scene_window_%s :\n", ident)
		fmt.Fprintf(&b, "    allActionsInSceneWindows sceneWindows timeOrder activeActions_%s scopeToScene := by\n", ident)
		fmt.Fprintf(&b, "  native_decide\n\n")
		fmt.Fprintf(&b, "theorem no_forbidden_states_%s :\n", ident)
		fmt.Fprintf(&b, "    noForbiddenStatesRegistered projectRules_%s stateDecls_%s := by\n", ident, ident)
		fmt.Fprintf(&b, "  native_decide\n\n")
		fmt.Fprintf(&b, "theorem no_forbidden_transitions_%s :\n", ident)
		fmt.Fprintf(&b, "    allActionsRespectRules projectRules_%s activeActions_%s := by\n", ident, ident)
		fmt.Fprintf(&b, "  native_decide\n\n")
		fmt.Fprintf(&b, "theorem fixed_facts_stable_%s :\n", ident)
		fmt.Fprintf(&b, "    fixedFactsStable fixedFacts_%s stateDecls_%s activeActions_%s timeOrder := by\n", ident, ident, ident)
		fmt.Fprintf(&b, "  native_decide\n\n")
	}
	for _, r := range d.Rules {
		if r.Kind == project.RuleForbidState && lastTime != "" {
			bid := leanIdent(project.NormalizeBranch(r.Branch))
			fmt.Fprintf(&b, "theorem forbid_state_%s_%s_%s_at_end :\n",
				bid, leanIdent(r.Thing), leanPred(r.Pred))
			fmt.Fprintf(&b, "    ¬ listContains (predsAt fixedFacts_%s stateDecls_%s activeActions_%s timeOrder TimeId.%s ThingId.%s) PredId.%s := by\n",
				bid, bid, bid, leanIdent(lastTime), leanIdent(r.Thing), leanPred(r.Pred))
			fmt.Fprintf(&b, "  native_decide\n\n")
		}
	}
	fmt.Fprintf(&b, "end %s\n", ns)
	return b.String()
}

func namespace(d *project.Data) string {
	title := strings.TrimSpace(d.Meta.Title)
	switch title {
	case "桃太郎":
		return "Momotaro"
	default:
		name := leanIdent(title)
		if isASCIIIdent(name) {
			return name
		}
		return "Work"
	}
}

func isASCIIIdent(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		if r > 127 {
			return false
		}
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '_' {
			return false
		}
	}
	return true
}

func genInductive(name string, ids []string) string {
	var b strings.Builder
	fmt.Fprintf(&b, "inductive %s\n", name)
	if len(ids) == 0 {
		fmt.Fprintf(&b, "  | none_\n")
	} else {
		for _, id := range ids {
			ctor := leanIdent(id)
			if name == "PredId" {
				ctor = leanPred(id)
			}
			fmt.Fprintf(&b, "  | %s\n", ctor)
		}
	}
	fmt.Fprintf(&b, "  deriving DecidableEq, Repr\n\n")
	return b.String()
}

func thingIDs(d *project.Data) []string {
	ids := make([]string, 0, len(d.Things))
	for _, t := range d.Things {
		ids = append(ids, t.ID)
	}
	sort.Strings(ids)
	return ids
}

func actionDeclExpr(a project.Action) string {
	from := "none"
	if a.From != "" {
		from = "some PredId." + leanPred(a.From)
	}
	return fmt.Sprintf("⟨ThingId.%s, %s, PredId.%s, TimeId.%s, %s⟩",
		leanIdent(a.Thing), from, leanPred(a.To), leanIdent(a.At), scopeExpr(a.Scope))
}

func branchIDs(d *project.Data) []string {
	return d.AllBranchIDs()
}

func sceneIDs(d *project.Data) []string {
	ids := make([]string, 0, len(d.Scenes))
	for _, s := range d.Scenes {
		ids = append(ids, s.ID)
	}
	sort.Strings(ids)
	return ids
}

func predIDs(d *project.Data) []string {
	m := d.Preds()
	ids := make([]string, 0, len(m))
	for id := range m {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	return ids
}

func joinIDs(ids []string, prefix string) string {
	parts := make([]string, len(ids))
	for i, id := range ids {
		parts[i] = prefix + "." + leanIdent(id)
	}
	return strings.Join(parts, ", ")
}

func scopeExpr(scope string) string {
	if scope == "" || scope == "plot" {
		return "Scope.plot"
	}
	if strings.HasPrefix(scope, "novel:") {
		sid := strings.TrimPrefix(scope, "novel:")
		return fmt.Sprintf("Scope.novel_%s", leanIdent(sid))
	}
	return "Scope.plot"
}

func leanIdent(id string) string {
	if id == "" {
		return "none_"
	}
	runes := []rune(id)
	var b strings.Builder
	for i, r := range runes {
		switch {
		case r == '-' || r == '.':
			b.WriteRune('_')
		case unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_':
			if i == 0 && unicode.IsDigit(r) {
				b.WriteRune('t')
			}
			b.WriteRune(r)
		default:
			b.WriteRune('_')
		}
	}
	out := b.String()
	if out == "" {
		return "id_"
	}
	return out
}

var predRomanization = map[string]string{
	"人間":    "ningen",
	"動物":    "doubutsu",
	"赤ちゃん":  "akachan",
	"青年":    "seinen",
	"村在住":   "murazaiju",
	"旅立ち":   "tabidachi",
	"野良":    "nora",
	"仲間":    "nakama",
	"健在":    "kenzen",
	"退治済み":  "taijizumi",
	"鬼退治済み": "onitaijizumi",
}

func leanPred(pred string) string {
	if pred == "" {
		return "none_"
	}
	if r, ok := predRomanization[pred]; ok {
		return r
	}
	out := leanIdent(pred)
	if isASCIIIdent(out) {
		return out
	}
	sum := sha256.Sum256([]byte(pred))
	return "pred_" + hex.EncodeToString(sum[:4])
}
