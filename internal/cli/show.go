package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/spf13/cobra"

	"novel-logic/internal/project"
)

const novelPreviewLen = 200

var timelineCmd = &cobra.Command{
	Use:   "timeline",
	Short: "Show time axis with scenes and actions",
	RunE:  runTimeline,
}

var sceneListCmd = &cobra.Command{
	Use:   "list",
	Short: "List scenes",
	RunE:  runSceneList,
}

var sceneShowCmd = &cobra.Command{
	Use:   "show <id>",
	Short: "Show one scene",
	Args:  cobra.ExactArgs(1),
	RunE:  runSceneShow,
}

var thingListCmd = &cobra.Command{
	Use:   "list",
	Short: "List things",
	RunE:  runThingList,
}

var thingShowCmd = &cobra.Command{
	Use:   "show <id>",
	Short: "Show one thing",
	Args:  cobra.ExactArgs(1),
	RunE:  runThingShow,
}

var timeListCmd = &cobra.Command{
	Use:   "list",
	Short: "List time points in order",
	RunE:  runTimeList,
}

var factListCmd = &cobra.Command{
	Use:   "list",
	Short: "List facts",
	RunE:  runFactList,
}

var factShowCmd = &cobra.Command{
	Use:   "show <id>",
	Short: "Show one fact",
	Args:  cobra.ExactArgs(1),
	RunE:  runFactShow,
}

var actionListCmd = &cobra.Command{
	Use:   "list",
	Short: "List actions",
	RunE:  runActionList,
}

var actionShowCmd = &cobra.Command{
	Use:   "show <id>",
	Short: "Show one action",
	Args:  cobra.ExactArgs(1),
	RunE:  runActionShow,
}

var ruleListCmd = &cobra.Command{
	Use:   "list",
	Short: "List rules",
	RunE:  runRuleList,
}

var ruleShowCmd = &cobra.Command{
	Use:   "show <id>",
	Short: "Show one rule",
	Args:  cobra.ExactArgs(1),
	RunE:  runRuleShow,
}

var novelListCmd = &cobra.Command{
	Use:   "list",
	Short: "List novels",
	RunE:  runNovelList,
}

var novelShowCmd = &cobra.Command{
	Use:   "show <scene_id>",
	Short: "Show novel for a scene",
	Args:  cobra.ExactArgs(1),
	RunE:  runNovelShow,
}

func runTimeline(cmd *cobra.Command, args []string) error {
	d, err := loadProject()
	if err != nil {
		return exitErr(4, err)
	}
	verbose, _ := cmd.Flags().GetBool("verbose")
	branch, _ := cmd.Flags().GetString("branch")
	branch = project.NormalizeBranch(branch)
	if branch != project.MainBranch || cmd.Flags().Changed("branch") {
		fmt.Printf("branch: %s\n\n", branch)
	}
	fmt.Println("time_order:")
	for _, t := range d.Meta.TimeOrder {
		markers := timelineMarkers(d, t)
		line := fmt.Sprintf("  %s", t)
		if markers != "" {
			line += "  " + markers
		}
		fmt.Println(line)
	}
	fmt.Println()
	fmt.Println("scenes:")
	for _, s := range d.Scenes {
		fmt.Printf("  %s [%s..%s] %s\n", s.ID, s.TimeStart, s.TimeEnd, s.Summary)
	}
	fmt.Println()
	fmt.Println("actions:")
	for _, a := range d.ActiveActions(branch) {
		from := a.From
		if from == "" {
			from = "∅"
		}
		b := a.Branch
		if b == "" {
			b = project.MainBranch
		}
		fmt.Printf("  %s @ %s [%s]  %s: %s → %s", a.ID, a.At, b, a.Thing, from, a.To)
		if a.Label != "" {
			fmt.Printf(" (%s)", a.Label)
		}
		fmt.Println()
	}
	if verbose {
		fmt.Println()
		fmt.Println("facts:")
		for _, f := range d.Facts {
			fmt.Printf("  %s [%s] %s は %s (scope=%s)\n", f.ID, f.Kind, f.Thing, f.Pred, scopeOrPlot(f.Scope))
		}
	}
	return nil
}

func timelineMarkers(d *project.Data, t string) string {
	var parts []string
	for _, s := range d.Scenes {
		if s.TimeStart == t {
			parts = append(parts, "scene:"+s.ID+" start")
		}
		if s.TimeEnd == t {
			parts = append(parts, "scene:"+s.ID+" end")
		}
	}
	for _, a := range d.Actions {
		if a.At == t {
			parts = append(parts, "action:"+a.ID)
		}
	}
	return strings.Join(parts, ", ")
}

func runSceneList(cmd *cobra.Command, args []string) error {
	d, err := loadProject()
	if err != nil {
		return exitErr(4, err)
	}
	for _, s := range d.Scenes {
		novel := novelStatus(d, s.ID)
		fmt.Printf("%s [%s..%s] %s  novel=%s\n", s.ID, s.TimeStart, s.TimeEnd, s.Summary, novel)
	}
	return nil
}

func runSceneShow(cmd *cobra.Command, args []string) error {
	d, err := loadProject()
	if err != nil {
		return exitErr(4, err)
	}
	id := args[0]
	for _, s := range d.Scenes {
		if s.ID != id {
			continue
		}
		fmt.Printf("id: %s\n", s.ID)
		fmt.Printf("summary: %s\n", s.Summary)
		fmt.Printf("time: %s .. %s\n", s.TimeStart, s.TimeEnd)
		fmt.Printf("novel: %s\n", novelStatus(d, s.ID))
		nf, na, pa := d.SceneLayerCounts(s.ID)
		fmt.Printf("layers: novel(facts=%d actions=%d) plot_inherited(actions=%d)\n", nf, na, pa)
		printSceneThingRefs(d, d.SceneThingRefs(s.ID))
		return nil
	}
	return exitErrf(4, "scene %q not found", id)
}

func sceneActions(d *project.Data, sceneID string) []project.Action {
	s := findScene(d, sceneID)
	if s == nil {
		return nil
	}
	var out []project.Action
	for _, a := range d.Actions {
		if d.TimeLE(s.TimeStart, a.At) && d.TimeLE(a.At, s.TimeEnd) {
			out = append(out, a)
		}
	}
	return out
}

func findScene(d *project.Data, id string) *project.Scene {
	for i := range d.Scenes {
		if d.Scenes[i].ID == id {
			return &d.Scenes[i]
		}
	}
	return nil
}

func novelStatus(d *project.Data, sceneID string) string {
	for _, n := range d.Novels {
		if n.SceneID == sceneID && project.NormalizeBranch(n.Branch) == project.MainBranch {
			return "yes"
		}
	}
	return "no"
}

func runThingList(cmd *cobra.Command, args []string) error {
	d, err := loadProject()
	if err != nil {
		return exitErr(4, err)
	}
	tags, _ := cmd.Flags().GetStringArray("tag")
	scope, _ := cmd.Flags().GetString("scope")
	printListFilters(tags, scope, "", "")
	n := 0
	for _, t := range d.Things {
		if !matchTags(t.Tags, tags) {
			continue
		}
		if scope != "" && scope != "all" && !containsStr(t.Scopes, scope) {
			continue
		}
		name := t.Name
		if name == "" {
			name = "-"
		}
		fixed, state := factCountsForThing(d, t.ID)
		fmt.Printf("%s name=%s tags=%v scopes=%v facts(fixed=%d state=%d)\n",
			t.ID, name, t.Tags, t.Scopes, fixed, state)
		n++
	}
	if n == 0 {
		fmt.Println("(no matches)")
	}
	return nil
}

func runThingShow(cmd *cobra.Command, args []string) error {
	d, err := loadProject()
	if err != nil {
		return exitErr(4, err)
	}
	t, _ := d.FindThing(args[0])
	if t == nil {
		return exitErrf(4, "thing %q not found", args[0])
	}
	fmt.Printf("id: %s\n", t.ID)
	if t.Name != "" {
		fmt.Printf("name: %s\n", t.Name)
	}
	fmt.Printf("tags: %v\n", t.Tags)
	fmt.Printf("scopes: %v\n", t.Scopes)
	fmt.Println("facts:")
	for _, f := range d.Facts {
		if f.Thing != t.ID {
			continue
		}
		fmt.Printf("  %s [%s] %s (scope=%s)\n", f.ID, f.Kind, f.Pred, scopeOrPlot(f.Scope))
	}
	fmt.Println("actions:")
	for _, a := range d.Actions {
		if a.Thing != t.ID {
			continue
		}
		fmt.Printf("  %s: %s → %s @ %s\n", a.ID, orEmpty(a.From), a.To, a.At)
	}
	return nil
}

func factCountsForThing(d *project.Data, thingID string) (fixed, state int) {
	for _, f := range d.Facts {
		if f.Thing != thingID {
			continue
		}
		switch f.Kind {
		case project.FactFixed:
			fixed++
		case project.FactState:
			state++
		}
	}
	return fixed, state
}

func runTimeList(cmd *cobra.Command, args []string) error {
	d, err := loadProject()
	if err != nil {
		return exitErr(4, err)
	}
	for i, t := range d.Meta.TimeOrder {
		fmt.Printf("%2d  %s\n", i+1, t)
	}
	return nil
}

func runFactList(cmd *cobra.Command, args []string) error {
	d, err := loadProject()
	if err != nil {
		return exitErr(4, err)
	}
	kind, _ := cmd.Flags().GetString("kind")
	thing, _ := cmd.Flags().GetString("thing")
	scope, _ := cmd.Flags().GetString("scope")
	tags, _ := cmd.Flags().GetStringArray("tag")
	printListFilters(tags, scope, kind, thing)
	n := 0
	for _, f := range d.Facts {
		if !matchKindFilter(kind, f.Kind) {
			continue
		}
		if thing != "" && f.Thing != thing {
			continue
		}
		if scope != "" && f.Scope != scope && !(scope == "plot" && f.Scope == "") {
			continue
		}
		if !thingMatchesTags(d, f.Thing, tags) {
			continue
		}
		fmt.Printf("%s [%s] %s は %s (scope=%s)\n", f.ID, f.Kind, f.Thing, f.Pred, scopeOrPlot(f.Scope))
		n++
	}
	if n == 0 {
		fmt.Println("(no matches)")
	}
	return nil
}

func runFactShow(cmd *cobra.Command, args []string) error {
	d, err := loadProject()
	if err != nil {
		return exitErr(4, err)
	}
	for _, f := range d.Facts {
		if f.ID != args[0] {
			continue
		}
		fmt.Printf("id: %s\n", f.ID)
		fmt.Printf("kind: %s\n", f.Kind)
		fmt.Printf("thing: %s\n", f.Thing)
		fmt.Printf("pred: %s\n", f.Pred)
		fmt.Printf("scope: %s\n", scopeOrPlot(f.Scope))
		return nil
	}
	return exitErrf(4, "fact %q not found", args[0])
}

func runActionList(cmd *cobra.Command, args []string) error {
	d, err := loadProject()
	if err != nil {
		return exitErr(4, err)
	}
	thing, _ := cmd.Flags().GetString("thing")
	tags, _ := cmd.Flags().GetStringArray("tag")
	printListFilters(tags, "", "", thing)
	n := 0
	for _, a := range d.Actions {
		if thing != "" && a.Thing != thing {
			continue
		}
		if !thingMatchesTags(d, a.Thing, tags) {
			continue
		}
		fmt.Printf("%s %s: %s → %s @ %s (scope=%s)", a.ID, a.Thing, orEmpty(a.From), a.To, a.At, scopeOrPlot(a.Scope))
		if a.Label != "" {
			fmt.Printf(" label=%q", a.Label)
		}
		fmt.Println()
		n++
	}
	if n == 0 {
		fmt.Println("(no matches)")
	}
	return nil
}

func runActionShow(cmd *cobra.Command, args []string) error {
	d, err := loadProject()
	if err != nil {
		return exitErr(4, err)
	}
	for _, a := range d.Actions {
		if a.ID != args[0] {
			continue
		}
		fmt.Printf("id: %s\n", a.ID)
		fmt.Printf("thing: %s\n", a.Thing)
		fmt.Printf("from: %s\n", orEmpty(a.From))
		fmt.Printf("to: %s\n", a.To)
		fmt.Printf("at: %s\n", a.At)
		fmt.Printf("scope: %s\n", scopeOrPlot(a.Scope))
		if a.Label != "" {
			fmt.Printf("label: %s\n", a.Label)
		}
		return nil
	}
	return exitErrf(4, "action %q not found", args[0])
}

func runRuleList(cmd *cobra.Command, args []string) error {
	d, err := loadProject()
	if err != nil {
		return exitErr(4, err)
	}
	for _, r := range d.Rules {
		switch r.Kind {
		case project.RuleForbidState:
			fmt.Printf("%s [forbid-state] %s は %s ではない\n", r.ID, r.Thing, r.Pred)
		case project.RuleForbidTransition:
			fmt.Printf("%s [forbid-transition] %s → %s は禁止\n", r.ID, r.From, r.To)
		default:
			fmt.Printf("%s [%s]\n", r.ID, r.Kind)
		}
	}
	return nil
}

func runRuleShow(cmd *cobra.Command, args []string) error {
	d, err := loadProject()
	if err != nil {
		return exitErr(4, err)
	}
	for _, r := range d.Rules {
		if r.ID != args[0] {
			continue
		}
		fmt.Printf("id: %s\n", r.ID)
		fmt.Printf("kind: %s\n", r.Kind)
		if r.Thing != "" {
			fmt.Printf("thing: %s\n", r.Thing)
		}
		if r.Pred != "" {
			fmt.Printf("pred: %s\n", r.Pred)
		}
		if r.From != "" {
			fmt.Printf("from: %s\n", r.From)
		}
		if r.To != "" {
			fmt.Printf("to: %s\n", r.To)
		}
		return nil
	}
	return exitErrf(4, "rule %q not found", args[0])
}

func runNovelList(cmd *cobra.Command, args []string) error {
	d, err := loadProject()
	if err != nil {
		return exitErr(4, err)
	}
	if len(d.Novels) == 0 {
		fmt.Println("(none)")
		return nil
	}
	for _, n := range d.Novels {
		rev := n.Revision
		if rev != "" && len(n.Revisions) > 0 && n.Revisions[len(n.Revisions)-1].Short != "" {
			rev = n.Revisions[len(n.Revisions)-1].Short
		}
		br := project.NormalizeBranch(n.Branch)
		if rev != "" {
			fmt.Printf("%s branch=%s [%s..%s] path=%s rev=%s\n", n.SceneID, br, n.TimeStart, n.TimeEnd, n.BodyPath, rev)
		} else {
			fmt.Printf("%s branch=%s [%s..%s] path=%s\n", n.SceneID, br, n.TimeStart, n.TimeEnd, n.BodyPath)
		}
	}
	for _, s := range d.Scenes {
		if novelStatus(d, s.ID) == "no" {
			fmt.Printf("%s [no novel]\n", s.ID)
		}
	}
	return nil
}

func runNovelShow(cmd *cobra.Command, args []string) error {
	d, err := loadProject()
	if err != nil {
		return exitErr(4, err)
	}
	full, _ := cmd.Flags().GetBool("full")
	branch, _ := cmd.Flags().GetString("branch")
	sceneID := args[0]
	meta, _ := d.FindNovel(sceneID, branch)
	if meta == nil {
		return exitErrf(4, "novel for scene %q on branch %q not found", sceneID, project.NormalizeBranch(branch))
	}
	body, err := os.ReadFile(filepath.Join(d.Root, meta.BodyPath))
	if err != nil {
		return exitErr(4, err)
	}
	text := string(body)
	nf, na, pa := d.SceneLayerCounts(sceneID)
	fmt.Printf("scene_id: %s\n", meta.SceneID)
	fmt.Printf("branch: %s\n", project.NormalizeBranch(meta.Branch))
	fmt.Printf("time: %s .. %s\n", meta.TimeStart, meta.TimeEnd)
	printNovelBodyGitInfo(d.Root, sceneID, meta)
	fmt.Printf("layers: novel(facts=%d actions=%d) plot_inherited(actions=%d)\n", nf, na, pa)
	printSceneThingRefs(d, d.SceneThingRefs(sceneID))
	printAlignmentNote(nf, na, pa)
	fmt.Println("body:")
	if full || len(text) <= novelPreviewLen {
		fmt.Println(text)
	} else {
		fmt.Println(text[:novelPreviewLen] + "…")
	}
	return nil
}

func printSceneThingRefs(d *project.Data, refs []project.ThingRef) {
	if len(refs) == 0 {
		fmt.Println("related_things: (none)")
		return
	}
	byLayer := map[string][]project.ThingRef{}
	for _, r := range refs {
		byLayer[r.Layer] = append(byLayer[r.Layer], r)
	}
	order := []string{"novel", "plot"}
	total := 0
	for _, layer := range order {
		total += len(byLayer[layer])
	}
	fmt.Printf("related_things (%d):\n", total)
	for _, layer := range order {
		items := byLayer[layer]
		if len(items) == 0 {
			fmt.Printf("  [%s] (none)\n", layer)
			continue
		}
		sort.Slice(items, func(i, j int) bool { return items[i].ID < items[j].ID })
		for _, r := range items {
			name := "-"
			tags := []string{}
			if t, _ := d.FindThing(r.ID); t != nil {
				if t.Name != "" {
					name = t.Name
				}
				tags = t.Tags
			}
			fmt.Printf("  [%s] %s name=%s tags=%v  (%s)\n", layer, r.ID, name, tags, r.Reason)
		}
	}
}

func printNovelBodyGitInfo(root, sceneID string, meta *project.NovelMeta) {
	fmt.Println("body_file:")
	fmt.Printf("  path: %s\n", meta.BodyPath)
	fmt.Println("  edit: open this file in your editor (novel-logic does not write prose)")
	fmt.Println("git:")
	if !project.IsGitRepo(root) {
		fmt.Println("  status: no git repository detected in this project")
		fmt.Println("  note: revision pinning is optional and requires git")
		return
	}
	if meta.Revision == "" {
		fmt.Println("  pinned_commit: (none)")
		fmt.Println("  workflow:")
		fmt.Printf("    1) edit %s\n", meta.BodyPath)
		fmt.Println("    2) git add / git commit")
		fmt.Printf("    3) novel-logic novel revision pin %s\n", sceneID)
		fmt.Println("  why: records which git commit this scene body matches (for check / CI)")
		return
	}
	fmt.Printf("  pinned_commit: %s\n", meta.Revision)
	if len(meta.Revisions) > 0 {
		r := meta.Revisions[len(meta.Revisions)-1]
		if r.Short != "" {
			fmt.Printf("  short: %s\n", r.Short)
		}
		if r.Branch != "" {
			fmt.Printf("  branch: %s\n", r.Branch)
		}
		if !r.RecordedAt.IsZero() {
			fmt.Printf("  pinned_at: %s\n", r.RecordedAt.UTC().Format("2006-01-02T15:04:05Z"))
		}
		if r.Note != "" {
			fmt.Printf("  note: %s\n", r.Note)
		}
		if r.Dirty {
			fmt.Println("  warning: pinned while file had uncommitted changes")
		}
	}
	fmt.Println("  update: after editing and committing, run novel revision pin again")
}

func printAlignmentNote(novelFacts, novelActions, plotActions int) {
	if novelFacts > 0 || novelActions > 0 {
		if plotActions > 0 {
			fmt.Println("alignment: mixed — novel layer + plot-inherited actions in scene window")
		} else {
			fmt.Println("alignment: novel layer only (scope matches scene)")
		}
		return
	}
	if plotActions > 0 {
		fmt.Println("alignment: plot-inherited only — re-register with --scope novel:<scene> for Phase B")
	}
}

func matchKindFilter(filter string, kind project.FactKind) bool {
	switch filter {
	case "", "all":
		return true
	case "fixed":
		return kind == project.FactFixed
	case "state":
		return kind == project.FactState
	default:
		return true
	}
}

func scopeOrPlot(s string) string {
	if s == "" {
		return "plot"
	}
	return s
}

func orEmpty(s string) string {
	if s == "" {
		return "∅"
	}
	return s
}

func containsStr(xs []string, x string) bool {
	for _, v := range xs {
		if v == x {
			return true
		}
	}
	return false
}

func matchTags(thingTags, filterTags []string) bool {
	if len(filterTags) == 0 {
		return true
	}
	for _, ft := range filterTags {
		if !containsStr(thingTags, ft) {
			return false
		}
	}
	return true
}

func thingMatchesTags(d *project.Data, thingID string, filterTags []string) bool {
	if len(filterTags) == 0 {
		return true
	}
	t, _ := d.FindThing(thingID)
	if t == nil {
		return false
	}
	return matchTags(t.Tags, filterTags)
}

func printListFilters(tags []string, scope, kind, thing string) {
	var parts []string
	if len(tags) > 0 {
		parts = append(parts, fmt.Sprintf("tag=%v", tags))
	}
	if scope != "" && scope != "all" {
		parts = append(parts, "scope="+scope)
	}
	if kind != "" && kind != "all" {
		parts = append(parts, "kind="+kind)
	}
	if thing != "" {
		parts = append(parts, "thing="+thing)
	}
	if len(parts) > 0 {
		fmt.Printf("# filter: %s\n", strings.Join(parts, ", "))
	}
}

func initShowCommands() {
	timelineCmd.Flags().Bool("verbose", false, "include all facts")
	timelineCmd.Flags().String("branch", project.MainBranch, "story branch for active actions")

	novelShowCmd.Flags().String("branch", project.MainBranch, "story branch id")

	thingListCmd.Flags().StringArray("tag", nil, "filter by tag (repeatable; thing must have all listed tags)")
	thingListCmd.Flags().String("scope", "all", "filter by scope (plot, novel:<scene>, all)")

	factListCmd.Flags().String("kind", "all", "fixed, state, or all")
	factListCmd.Flags().String("thing", "", "filter by thing id")
	factListCmd.Flags().String("scope", "", "filter by scope")
	factListCmd.Flags().StringArray("tag", nil, "filter by subject thing tag (repeatable; AND)")

	actionListCmd.Flags().String("thing", "", "filter by thing id")
	actionListCmd.Flags().StringArray("tag", nil, "filter by subject thing tag (repeatable; AND)")

	novelShowCmd.Flags().Bool("full", false, "print full body text")

	sceneCmd.AddCommand(sceneListCmd, sceneShowCmd)
	thingCmd.AddCommand(thingListCmd, thingShowCmd)
	timeCmd.AddCommand(timeListCmd)
	factCmd.AddCommand(factListCmd, factShowCmd)
	actionCmd.AddCommand(actionListCmd, actionShowCmd)
	ruleCmd.AddCommand(ruleListCmd, ruleShowCmd)
	novelCmd.AddCommand(novelListCmd, novelShowCmd)
	rootCmd.AddCommand(timelineCmd)
}