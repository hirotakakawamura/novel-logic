package project

import "fmt"

// TimeRegistryIssues reports mismatches between times.yaml and project.yaml time_order.
func TimeRegistryIssues(d *Data) []string {
	registry := make(map[string]bool)
	for _, te := range d.Times {
		if te.ID == "" {
			continue
		}
		registry[te.ID] = true
	}
	seenOrder := make(map[string]bool)
	var issues []string
	for _, id := range d.Meta.TimeOrder {
		if id == "" {
			issues = append(issues, "empty time id in time_order")
			continue
		}
		if seenOrder[id] {
			issues = append(issues, fmt.Sprintf("duplicate time %q in time_order", id))
			continue
		}
		seenOrder[id] = true
		if !registry[id] {
			issues = append(issues, fmt.Sprintf("time %q in time_order but missing from times.yaml", id))
		}
	}
	for id := range registry {
		if !seenOrder[id] {
			issues = append(issues, fmt.Sprintf("time %q in times.yaml but missing from time_order", id))
		}
	}
	return issues
}