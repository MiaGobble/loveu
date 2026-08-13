package runtime_test

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	rt "github.com/MiaGobble/loveu/cli/internal/runtime"
)

func TestProbeVersion(t *testing.T) {
	dir := t.TempDir()
	var script string
	var exe string
	if runtime.GOOS == "windows" {
		exe = filepath.Join(dir, "love.bat")
		script = "@echo off\r\necho loveu 0.1.0 (LOVE 12.0 \"Bestest Friend\")\r\n"
	} else {
		exe = filepath.Join(dir, "love")
		script = "#!/bin/sh\necho 'loveu 0.1.0 (LÖVE 12.0 \"Bestest Friend\")'\n"
	}
	if err := os.WriteFile(exe, []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}
	got, _, err := rt.ProbeVersion(exe)
	if err != nil {
		t.Fatal(err)
	}
	if got != "0.1.0" {
		t.Fatalf("got %q", got)
	}
}
