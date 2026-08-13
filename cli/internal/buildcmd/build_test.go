package buildcmd_test

import (
	"archive/zip"
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/MiaGobble/loveu/cli/internal/buildcmd"
)

func TestWriteLoveArchive(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "loveu.toml"), []byte("name=\"t\"\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(root, "src"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "src", "main.luau"), []byte("print(1)\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(root, "dist"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "dist", "skip.txt"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	out := filepath.Join(t.TempDir(), "game.love")
	if err := buildcmd.WriteLoveArchive(root, out); err != nil {
		t.Fatal(err)
	}
	r, err := zip.OpenReader(out)
	if err != nil {
		t.Fatal(err)
	}
	defer r.Close()
	names := map[string]bool{}
	for _, f := range r.File {
		names[f.Name] = true
	}
	if !names["loveu.toml"] || !names["src/main.luau"] {
		t.Fatalf("missing entries: %v", names)
	}
	for n := range names {
		if n == "dist/skip.txt" || n == "dist/" {
			t.Fatalf("dist should be excluded, found %s", n)
		}
	}
}

func TestFuseBinary(t *testing.T) {
	dir := t.TempDir()
	exe := filepath.Join(dir, "love.exe")
	love := filepath.Join(dir, "g.love")
	dest := filepath.Join(dir, "out.exe")
	if err := os.WriteFile(exe, []byte("EXE"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(love, []byte("LOVE"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := buildcmd.FuseBinary(exe, love, dest); err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile(dest)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, []byte("EXELOVE")) {
		t.Fatalf("got %q", got)
	}
}
