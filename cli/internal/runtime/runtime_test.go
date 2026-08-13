package runtime_test

import (
	"archive/zip"
	"os"
	"path/filepath"
	"testing"

	"github.com/MiaGobble/loveu/cli/internal/runtime"
)

func TestAssetName(t *testing.T) {
	cases := map[runtime.Platform]string{
		runtime.WindowsAMD64: "loveu-engine-0.1.0-windows-x86_64.zip",
		runtime.MacOS:        "loveu-engine-0.1.0-macos.zip",
		runtime.LinuxAMD64:   "loveu-engine-0.1.0-linux-x86_64.AppImage",
		runtime.Android:      "loveu-engine-0.1.0-android.apk",
		runtime.IOS:          "loveu-engine-0.1.0-ios.zip",
	}
	for p, want := range cases {
		if got := runtime.AssetName("0.1.0", p); got != want {
			t.Errorf("%s: got %s want %s", p, got, want)
		}
	}
}

func TestInstallFromLocalZip(t *testing.T) {
	home := t.TempDir()
	t.Setenv("LOVEU_HOME", home)

	dir := t.TempDir()
	zipPath := filepath.Join(dir, "engine.zip")
	if err := writeMinimalZip(zipPath, "love.exe", []byte("fake")); err != nil {
		t.Fatal(err)
	}
	dest, err := runtime.Install(runtime.InstallOptions{
		Version:  "0.1.0",
		Platform: runtime.WindowsAMD64,
		From:     zipPath,
	})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(dest, "love.exe")); err != nil {
		t.Fatal(err)
	}
	exe, ok := runtime.Installed("0.1.0", runtime.WindowsAMD64)
	if !ok {
		t.Fatal("expected installed")
	}
	if filepath.Base(exe) != "love.exe" && filepath.Base(exe) != "windows-x86_64" {
		// Installed returns the cache directory.
		if filepath.Base(exe) != "windows-x86_64" {
			t.Fatalf("got %s", exe)
		}
	}
	bin, err := runtime.FindBinary(exe, runtime.WindowsAMD64, false)
	if err != nil {
		t.Fatal(err)
	}
	if filepath.Base(bin) != "love.exe" {
		t.Fatalf("got %s", bin)
	}
}

func TestParsePlatform(t *testing.T) {
	if _, err := runtime.ParsePlatform("nope"); err == nil {
		t.Fatal("expected error")
	}
	p, err := runtime.ParsePlatform("macos")
	if err != nil || p != runtime.MacOS {
		t.Fatalf("got %v %v", p, err)
	}
}

func writeMinimalZip(path, name string, data []byte) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	zw := zip.NewWriter(f)
	w, err := zw.Create(name)
	if err != nil {
		return err
	}
	if _, err := w.Write(data); err != nil {
		return err
	}
	return zw.Close()
}
