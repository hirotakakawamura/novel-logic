package lean

import (
	"fmt"
	"os/exec"
	"strings"
)

type Toolchain struct {
	Elan  string
	Lean  string
	Lake  string
	Found bool
}

func Detect() Toolchain {
	tc := Toolchain{}
	tc.Elan, _ = exec.LookPath("elan")
	tc.Lean, _ = exec.LookPath("lean")
	tc.Lake, _ = exec.LookPath("lake")
	tc.Found = tc.Lake != "" && tc.Lean != ""
	return tc
}

func (tc Toolchain) Version() string {
	if tc.Lean == "" {
		return "not found"
	}
	out, err := exec.Command(tc.Lean, "--version").CombinedOutput()
	if err != nil {
		return strings.TrimSpace(string(out))
	}
	return strings.TrimSpace(string(out))
}

func LakeBuild(logicDir string) (string, error) {
	tc := Detect()
	if !tc.Found {
		return "", fmt.Errorf("lean/lake not found in PATH")
	}
	cmd := exec.Command(tc.Lake, "build")
	cmd.Dir = logicDir
	out, err := cmd.CombinedOutput()
	return string(out), err
}