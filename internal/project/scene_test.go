package project

import "testing"

func TestParseNovelScope(t *testing.T) {
	tests := []struct {
		scope string
		want  string
	}{
		{"novel:scene1", "scene1"},
		{"plot", ""},
		{"novel:", ""},
		{"", ""},
	}
	for _, tt := range tests {
		if got := ParseNovelScope(tt.scope); got != tt.want {
			t.Errorf("ParseNovelScope(%q) = %q, want %q", tt.scope, got, tt.want)
		}
	}
}

func TestSceneThingRefsDualLayer(t *testing.T) {
	d := newTestProject(t)
	d.Things[0].Scopes = append(d.Things[0].Scopes, NovelScope("scene1"))
	d.Facts = append(d.Facts, Fact{
		ID: "nf1", Kind: FactState, Thing: "ally", Pred: "present",
		Scope: NovelScope("scene1"), Branch: MainBranch,
	})

	refs := d.SceneThingRefs("scene1")
	if len(refs) < 2 {
		t.Fatalf("refs = %+v, want hero (plot+novel) and ally (novel)", refs)
	}
	byID := map[string][]ThingRef{}
	for _, r := range refs {
		byID[r.ID] = append(byID[r.ID], r)
	}
	if _, ok := byID["hero"]; !ok {
		t.Fatalf("missing hero in %+v", refs)
	}
	if _, ok := byID["ally"]; !ok {
		t.Fatalf("missing ally in %+v", refs)
	}
	hasPlotHero := false
	hasNovelHero := false
	for _, r := range byID["hero"] {
		switch r.Layer {
		case "plot":
			hasPlotHero = true
		case "novel":
			hasNovelHero = true
		}
	}
	if !hasPlotHero {
		t.Fatal("hero missing plot layer ref (act1 at t2 in scene1 window)")
	}
	if !hasNovelHero {
		t.Fatal("hero missing novel layer ref (scope on thing)")
	}
}

func TestSceneThingRefsUnknownScene(t *testing.T) {
	d := newTestProject(t)
	if refs := d.SceneThingRefs("missing"); refs != nil {
		t.Fatalf("refs = %+v, want nil", refs)
	}
}

func TestSceneLayerCounts(t *testing.T) {
	d := newTestProject(t)
	d.Facts = append(d.Facts, Fact{
		ID: "nf1", Kind: FactState, Thing: "ally", Pred: "present",
		Scope: NovelScope("scene1"), Branch: MainBranch,
	})
	d.Actions = append(d.Actions, Action{
		ID: "na1", Thing: "hero", From: "mid", To: "end",
		At: "t1", Scope: NovelScope("scene1"), Branch: MainBranch,
	})

	nf, na, pa := d.SceneLayerCounts("scene1")
	if nf != 1 {
		t.Fatalf("novelFacts = %d, want 1", nf)
	}
	if na != 1 {
		t.Fatalf("novelActions = %d, want 1", na)
	}
	if pa != 1 {
		t.Fatalf("plotActions = %d, want 1 (act1 at t2)", pa)
	}
}

func TestSceneLayerCountsUnknownScene(t *testing.T) {
	d := newTestProject(t)
	nf, na, pa := d.SceneLayerCounts("ghost")
	if nf != 0 || na != 0 || pa != 0 {
		t.Fatalf("got (%d,%d,%d), want (0,0,0)", nf, na, pa)
	}
}