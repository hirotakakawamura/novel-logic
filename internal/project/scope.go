package project

// EnsureThingNovelScope adds novel:<sceneID> to thing scopes when registering novel-layer data.
func (d *Data) EnsureThingNovelScope(thingID, scope string) {
	sceneID := ParseNovelScope(scope)
	if sceneID == "" {
		return
	}
	t, _ := d.FindThing(thingID)
	if t == nil {
		return
	}
	t.Scopes = MergeScope(t.Scopes, NovelScope(sceneID))
}