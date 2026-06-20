package generate

import (
	"testing"

	"novel-logic/internal/project"
	"novel-logic/internal/testfixture"
)

func minimalProject(t *testing.T) *project.Data {
	return testfixture.LoadMinimal(t)
}