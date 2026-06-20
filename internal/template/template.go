package template

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

//go:embed all:data
var dataFS embed.FS

func Materialize(dest, name string) error {
	if name == "" {
		name = "default"
	}
	src := filepath.Join("data", name)
	info, err := fs.Stat(dataFS, src)
	if err != nil || !info.IsDir() {
		return fmt.Errorf("unknown template %q", name)
	}
	return fs.WalkDir(dataFS, src, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		if rel == "." {
			return nil
		}
		target := filepath.Join(dest, rel)
		if d.IsDir() {
			return os.MkdirAll(target, 0o755)
		}
		b, err := dataFS.ReadFile(path)
		if err != nil {
			return err
		}
		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
			return err
		}
		return os.WriteFile(target, b, 0o644)
	})
}

func List() ([]string, error) {
	entries, err := fs.ReadDir(dataFS, "data")
	if err != nil {
		return nil, err
	}
	var names []string
	for _, e := range entries {
		if e.IsDir() {
			names = append(names, e.Name())
		}
	}
	return names, nil
}