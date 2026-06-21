package project

import (
	"strings"
	"testing"
)

func TestTimeRegistryIssues(t *testing.T) {
	tests := []struct {
		name     string
		mutate   func(*Data)
		wantSubs []string
	}{
		{
			name:   "clean",
			mutate: func(*Data) {},
		},
		{
			name: "time_order_only",
			mutate: func(d *Data) {
				d.Meta.TimeOrder = append(d.Meta.TimeOrder, "ghost")
			},
			wantSubs: []string{`time "ghost" in time_order but missing from times.yaml`},
		},
		{
			name: "times_yaml_only",
			mutate: func(d *Data) {
				d.Times = append(d.Times, TimeEntry{ID: "orphan"})
			},
			wantSubs: []string{`time "orphan" in times.yaml but missing from time_order`},
		},
		{
			name: "duplicate_in_time_order",
			mutate: func(d *Data) {
				d.Meta.TimeOrder = append(d.Meta.TimeOrder, "t1")
			},
			wantSubs: []string{`duplicate time "t1" in time_order`},
		},
		{
			name: "empty_id_in_time_order",
			mutate: func(d *Data) {
				d.Meta.TimeOrder = append(d.Meta.TimeOrder, "")
			},
			wantSubs: []string{"empty time id in time_order"},
		},
		{
			name: "empty_id_in_times_yaml",
			mutate: func(d *Data) {
				d.Times = append(d.Times, TimeEntry{ID: ""})
			},
			wantSubs: []string{"empty time id in times.yaml"},
		},
		{
			name: "deterministic_orphan_order",
			mutate: func(d *Data) {
				d.Times = append(d.Times,
					TimeEntry{ID: "z_orphan"},
					TimeEntry{ID: "a_orphan"},
				)
			},
			wantSubs: []string{
				`time "a_orphan" in times.yaml but missing from time_order`,
				`time "z_orphan" in times.yaml but missing from time_order`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := newTestProject(t)
			tt.mutate(d)
			issues := TimeRegistryIssues(d)
			if len(tt.wantSubs) == 0 {
				if len(issues) != 0 {
					t.Fatalf("expected no issues, got %v", issues)
				}
				return
			}
			if len(issues) < len(tt.wantSubs) {
				t.Fatalf("got %d issues %v, want at least %d", len(issues), issues, len(tt.wantSubs))
			}
			joined := strings.Join(issues, "\n")
			for _, sub := range tt.wantSubs {
				if !strings.Contains(joined, sub) {
					t.Fatalf("issues %v missing %q", issues, sub)
				}
			}
			if tt.name == "deterministic_orphan_order" {
				idxA := strings.Index(joined, `time "a_orphan"`)
				idxZ := strings.Index(joined, `time "z_orphan"`)
				if idxA < 0 || idxZ < 0 || idxA > idxZ {
					t.Fatalf("expected a_orphan before z_orphan, got %v", issues)
				}
			}
		})
	}
}