package project

import "testing"

func TestTimeRegistryIssues(t *testing.T) {
	d := newTestProject(t)
	if issues := TimeRegistryIssues(d); len(issues) != 0 {
		t.Fatalf("expected no issues, got %v", issues)
	}

	d.Meta.TimeOrder = append(d.Meta.TimeOrder, "ghost")
	if issues := TimeRegistryIssues(d); len(issues) == 0 {
		t.Fatal("expected mismatch for time_order-only id")
	}

	d = newTestProject(t)
	d.Times = append(d.Times, TimeEntry{ID: "orphan"})
	if issues := TimeRegistryIssues(d); len(issues) == 0 {
		t.Fatal("expected mismatch for times.yaml-only id")
	}
}