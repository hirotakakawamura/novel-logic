package project

import "strings"

// NovelScope returns the scope string for a scene's novel layer.
func NovelScope(sceneID string) string {
	return "novel:" + sceneID
}

// ParseNovelScope returns sceneID if scope is novel:<id>, else "".
func ParseNovelScope(scope string) string {
	if strings.HasPrefix(scope, "novel:") {
		return strings.TrimPrefix(scope, "novel:")
	}
	return ""
}

// SceneWindow returns the time window for a scene.
func (d *Data) SceneWindow(sceneID string) (start, end string, ok bool) {
	for _, s := range d.Scenes {
		if s.ID == sceneID {
			return s.TimeStart, s.TimeEnd, true
		}
	}
	return "", "", false
}

// ScenesContainingTime returns scene IDs whose window contains time at.
func (d *Data) ScenesContainingTime(at string) []string {
	var ids []string
	for _, s := range d.Scenes {
		if d.TimeLE(s.TimeStart, at) && d.TimeLE(at, s.TimeEnd) {
			ids = append(ids, s.ID)
		}
	}
	return ids
}

// ThingRef is a thing referenced from a scene layer with provenance.
type ThingRef struct {
	ID     string
	Layer  string // "novel" or "plot"
	Reason string
}

// SceneThingRefs returns related things for a scene using the dual-layer model.
// novel layer: scope=novel:sceneID on things, facts, actions.
// plot layer: plot-scoped actions whose at falls in the scene time window.
func (d *Data) SceneThingRefs(sceneID string) []ThingRef {
	start, end, ok := d.SceneWindow(sceneID)
	if !ok {
		return nil
	}
	novelScope := NovelScope(sceneID)
	seen := map[string]map[string]bool{}
	var out []ThingRef

	add := func(id, layer, reason string) {
		if id == "" {
			return
		}
		if seen[id] == nil {
			seen[id] = map[string]bool{}
		}
		key := layer + "|" + reason
		if seen[id][key] {
			return
		}
		seen[id][key] = true
		out = append(out, ThingRef{ID: id, Layer: layer, Reason: reason})
	}

	for _, t := range d.Things {
		for _, sc := range t.Scopes {
			if sc == novelScope {
				add(t.ID, "novel", "scope:"+novelScope)
			}
		}
	}
	for _, f := range d.Facts {
		if f.Scope == novelScope {
			add(f.Thing, "novel", "fact:"+f.ID)
		}
	}
	for _, a := range d.Actions {
		scope := a.Scope
		if scope == "" {
			scope = "plot"
		}
		if scope == novelScope {
			add(a.Thing, "novel", "action:"+a.ID)
			continue
		}
		if scope == "plot" && d.TimeLE(start, a.At) && d.TimeLE(a.At, end) {
			add(a.Thing, "plot", "@"+a.At+":"+a.ID)
		}
	}
	return out
}

// SceneLayerCounts returns fact/action counts per layer for a scene.
func (d *Data) SceneLayerCounts(sceneID string) (novelFacts, novelActions, plotActions int) {
	start, end, ok := d.SceneWindow(sceneID)
	if !ok {
		return 0, 0, 0
	}
	novelScope := NovelScope(sceneID)
	for _, f := range d.Facts {
		if f.Scope == novelScope {
			novelFacts++
		}
	}
	for _, a := range d.Actions {
		scope := a.Scope
		if scope == "" {
			scope = "plot"
		}
		if scope == novelScope {
			novelActions++
			continue
		}
		if scope == "plot" && d.TimeLE(start, a.At) && d.TimeLE(a.At, end) {
			plotActions++
		}
	}
	return novelFacts, novelActions, plotActions
}