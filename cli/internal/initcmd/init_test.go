package initcmd_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/MiaGobble/loveu/cli/internal/initcmd"
	"github.com/MiaGobble/loveu/cli/internal/project"
)

func TestInitLayout(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "game")
	if err := initcmd.Run(initcmd.Options{Dir: dir, Name: "hello"}); err != nil {
		t.Fatal(err)
	}
	proj, err := project.Load(filepath.Join(dir, "loveu.toml"))
	if err != nil {
		t.Fatal(err)
	}
	if proj.Name != "hello" || proj.CodeRoot != "src" {
		t.Fatalf("%+v", proj)
	}
	if _, err := os.Stat(filepath.Join(dir, "src", "main.luau")); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(dir, ".gitignore")); err != nil {
		t.Fatal(err)
	}
}

func TestInitRefuseNonEmpty(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "x"), []byte("1"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := initcmd.Run(initcmd.Options{Dir: dir}); err == nil {
		t.Fatal("expected error")
	}
}
