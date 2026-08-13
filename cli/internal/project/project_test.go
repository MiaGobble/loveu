package project_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/MiaGobble/loveu/cli/internal/project"
)

func TestFindAndLoad(t *testing.T) {
	root := t.TempDir()
	nested := filepath.Join(root, "a", "b")
	if err := os.MkdirAll(nested, 0o755); err != nil {
		t.Fatal(err)
	}
	toml := `name = "demo"
version = "1.2.3"
engine_version = "0.1.0"
code_root = "src"
`
	if err := os.WriteFile(filepath.Join(root, "loveu.toml"), []byte(toml), 0o644); err != nil {
		t.Fatal(err)
	}
	proj, err := project.FindAndLoad(nested)
	if err != nil {
		t.Fatal(err)
	}
	if proj.Name != "demo" || proj.EngineVersion != "0.1.0" || proj.CodeRoot != "src" {
		t.Fatalf("unexpected project: %+v", proj)
	}
}

func TestMissingRequired(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "loveu.toml"), []byte(`name = "x"`+"\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	_, err := project.Load(filepath.Join(dir, "loveu.toml"))
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestNotFound(t *testing.T) {
	dir := t.TempDir()
	_, err := project.FindAndLoad(dir)
	if err != project.ErrNotFound {
		t.Fatalf("got %v", err)
	}
}
