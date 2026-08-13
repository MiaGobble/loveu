package buildcmd

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/MiaGobble/loveu/cli/internal/project"
	rt "github.com/MiaGobble/loveu/cli/internal/runtime"
)

type Options struct {
	Target   string
	Offline  bool
	Keystore string
}

var skipNames = map[string]bool{
	".git": true, "dist": true, ".DS_Store": true, "node_modules": true,
}

func Run(opts Options) error {
	proj, err := project.FindAndLoad(".")
	if err != nil {
		return err
	}
	target := strings.ToLower(opts.Target)
	distRoot := filepath.Join(proj.Root, "dist")
	if err := os.MkdirAll(distRoot, 0o755); err != nil {
		return err
	}

	lovePath := filepath.Join(distRoot, fmt.Sprintf("%s-%s.love", proj.Name, proj.Version))
	if err := WriteLoveArchive(proj.Root, lovePath); err != nil {
		return err
	}
	fmt.Printf("wrote %s\n", lovePath)

	switch target {
	case "love":
		return nil
	case "windows":
		return buildWindows(proj, lovePath, distRoot, opts.Offline)
	case "macos":
		return buildMacOS(proj, lovePath, distRoot, opts.Offline)
	case "linux":
		return buildLinux(proj, lovePath, distRoot, opts.Offline)
	case "android":
		return buildAndroid(proj, lovePath, distRoot, opts)
	case "ios":
		return buildIOS(proj, lovePath, distRoot, opts.Offline)
	default:
		return fmt.Errorf("unknown build target %q (want love, windows, macos, linux, android, ios)", opts.Target)
	}
}

// WriteLoveArchive zips the project root into a .love file.
func WriteLoveArchive(root, dest string) error {
	if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
		return err
	}
	f, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer f.Close()
	zw := zip.NewWriter(f)
	defer zw.Close()

	return filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		if rel == "." {
			return nil
		}
		base := filepath.Base(path)
		if skipNames[base] {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		name := filepath.ToSlash(rel)
		if d.IsDir() {
			_, err := zw.Create(name + "/")
			return err
		}
		w, err := zw.Create(name)
		if err != nil {
			return err
		}
		in, err := os.Open(path)
		if err != nil {
			return err
		}
		_, copyErr := io.Copy(w, in)
		in.Close()
		return copyErr
	})
}

func buildWindows(proj *project.Project, lovePath, distRoot string, offline bool) error {
	plat := rt.WindowsAMD64
	dir, err := rt.EnsureInstalled(rt.InstallOptions{Version: proj.EngineVersion, Platform: plat, Offline: offline})
	if err != nil {
		return err
	}
	loveExe, err := rt.FindBinary(dir, plat, false)
	if err != nil {
		return err
	}
	outDir := filepath.Join(distRoot, fmt.Sprintf("%s-%s-windows", proj.Name, proj.Version))
	_ = os.RemoveAll(outDir)
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return err
	}
	// Copy runtime DLLs / companion files from love.exe directory
	rtDir := filepath.Dir(loveExe)
	if err := copyDirContents(rtDir, outDir); err != nil {
		return err
	}
	fused := filepath.Join(outDir, proj.Name+".exe")
	if err := FuseBinary(loveExe, lovePath, fused); err != nil {
		return err
	}
	// Remove unfused love.exe if still present under that name
	_ = os.Remove(filepath.Join(outDir, "love.exe"))
	zipOut := outDir + ".zip"
	if err := zipDir(outDir, zipOut); err != nil {
		return err
	}
	fmt.Printf("wrote %s\n", zipOut)
	return nil
}

func buildMacOS(proj *project.Project, lovePath, distRoot string, offline bool) error {
	dir, err := rt.EnsureInstalled(rt.InstallOptions{Version: proj.EngineVersion, Platform: rt.MacOS, Offline: offline})
	if err != nil {
		return err
	}
	app := findLoveApp(dir)
	if app == "" {
		return fmt.Errorf("love.app not found in %s", dir)
	}
	outDir := filepath.Join(distRoot, fmt.Sprintf("%s-%s-macos", proj.Name, proj.Version))
	_ = os.RemoveAll(outDir)
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return err
	}
	destApp := filepath.Join(outDir, proj.Name+".app")
	if err := copyTree(app, destApp); err != nil {
		return err
	}
	res := filepath.Join(destApp, "Contents", "Resources")
	if err := os.MkdirAll(res, 0o755); err != nil {
		return err
	}
	if err := copyFile(lovePath, filepath.Join(res, proj.Name+".love")); err != nil {
		return err
	}
	plist := filepath.Join(destApp, "Contents", "Info.plist")
	_ = patchPlist(plist, proj.Name)
	zipOut := outDir + ".zip"
	if err := zipDir(outDir, zipOut); err != nil {
		return err
	}
	fmt.Printf("wrote %s\n", zipOut)
	return nil
}

func buildLinux(proj *project.Project, lovePath, distRoot string, offline bool) error {
	plat := rt.LinuxAMD64
	if runtime.GOARCH == "arm64" {
		plat = rt.LinuxARM64
	}
	dir, err := rt.EnsureInstalled(rt.InstallOptions{Version: proj.EngineVersion, Platform: plat, Offline: offline})
	if err != nil {
		return err
	}
	outDir := filepath.Join(distRoot, fmt.Sprintf("%s-%s-linux", proj.Name, proj.Version))
	_ = os.RemoveAll(outDir)
	if err := copyTree(dir, outDir); err != nil {
		return err
	}
	bin, err := rt.FindBinary(outDir, plat, false)
	if err != nil {
		return err
	}
	// Prefer fusing bin/love when present
	fused := bin
	if filepath.Base(bin) == "love" || strings.HasSuffix(bin, string(filepath.Separator)+"love") {
		fused = filepath.Join(filepath.Dir(bin), proj.Name)
	} else {
		fused = filepath.Join(outDir, proj.Name)
	}
	if err := FuseBinary(bin, lovePath, fused); err != nil {
		return err
	}
	_ = os.Chmod(fused, 0o755)
	if fused != bin {
		_ = os.Remove(bin)
	}
	tarOut := outDir + ".tar.gz"
	if err := tarGzDir(outDir, tarOut); err != nil {
		return err
	}
	fmt.Printf("wrote %s\n", tarOut)
	return nil
}

func buildAndroid(proj *project.Project, lovePath, distRoot string, opts Options) error {
	dir, err := rt.EnsureInstalled(rt.InstallOptions{Version: proj.EngineVersion, Platform: rt.Android, Offline: opts.Offline})
	if err != nil {
		return err
	}
	apk, err := rt.FindBinary(dir, rt.Android, false)
	if err != nil {
		return err
	}
	outAPK := filepath.Join(distRoot, fmt.Sprintf("%s-%s-android.apk", proj.Name, proj.Version))
	if err := copyFile(apk, outAPK); err != nil {
		return err
	}
	if err := replaceZipEntry(outAPK, "assets/game.love", lovePath); err != nil {
		return fmt.Errorf("inject game.love into apk: %w", err)
	}
	fmt.Printf("wrote %s (unsigned or as-downloaded)\n", outAPK)
	if _, err := exec.LookPath("zipalign"); err == nil {
		aligned := outAPK + ".aligned"
		cmd := exec.Command("zipalign", "-f", "4", outAPK, aligned)
		if err := cmd.Run(); err == nil {
			_ = os.Rename(aligned, outAPK)
		}
	}
	if opts.Keystore != "" {
		if _, err := exec.LookPath("apksigner"); err == nil {
			cmd := exec.Command("apksigner", "sign", "--ks", opts.Keystore, outAPK)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			_ = cmd.Run()
		}
	}
	return nil
}

func buildIOS(proj *project.Project, lovePath, distRoot string, offline bool) error {
	if runtime.GOOS != "darwin" {
		return fmt.Errorf("iOS builds require macOS and Xcode")
	}
	dir, err := rt.EnsureInstalled(rt.InstallOptions{Version: proj.EngineVersion, Platform: rt.IOS, Offline: offline})
	if err != nil {
		return err
	}
	outDir := filepath.Join(distRoot, fmt.Sprintf("%s-%s-ios", proj.Name, proj.Version))
	_ = os.RemoveAll(outDir)
	if err := copyTree(dir, outDir); err != nil {
		return err
	}
	// Drop .love into first Resources folder found
	var resDir string
	_ = filepath.WalkDir(outDir, func(path string, d os.DirEntry, err error) error {
		if err != nil || !d.IsDir() {
			return nil
		}
		if d.Name() == "Resources" {
			resDir = path
			return filepath.SkipAll
		}
		return nil
	})
	if resDir == "" {
		resDir = filepath.Join(outDir, "Resources")
		_ = os.MkdirAll(resDir, 0o755)
	}
	if err := copyFile(lovePath, filepath.Join(resDir, proj.Name+".love")); err != nil {
		return err
	}
	fmt.Printf("wrote %s (sign with Xcode / CODE_SIGN_IDENTITY as needed)\n", outDir)
	return nil
}

// FuseBinary concatenates exeBytes + loveBytes into dest.
func FuseBinary(exePath, lovePath, dest string) error {
	exe, err := os.ReadFile(exePath)
	if err != nil {
		return err
	}
	love, err := os.ReadFile(lovePath)
	if err != nil {
		return err
	}
	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()
	if _, err := out.Write(exe); err != nil {
		return err
	}
	_, err = out.Write(love)
	return err
}

func findLoveApp(dir string) string {
	var found string
	_ = filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.IsDir() && d.Name() == "love.app" {
			found = path
			return filepath.SkipAll
		}
		return nil
	})
	return found
}

func patchPlist(path, name string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	s := string(data)
	s = strings.ReplaceAll(s, "love.app", name+".app")
	s = strings.ReplaceAll(s, ">LÖVE<", ">"+name+"<")
	s = strings.ReplaceAll(s, ">Love<", ">"+name+"<")
	return os.WriteFile(path, []byte(s), 0o644)
}

func replaceZipEntry(zipPath, entry, filePath string) error {
	tmp := zipPath + ".tmp"
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return err
	}
	defer r.Close()
	out, err := os.Create(tmp)
	if err != nil {
		return err
	}
	zw := zip.NewWriter(out)
	replaced := false
	for _, f := range r.File {
		if f.Name == entry {
			payload, err := os.ReadFile(filePath)
			if err != nil {
				zw.Close()
				out.Close()
				_ = os.Remove(tmp)
				return err
			}
			w, err := zw.Create(entry)
			if err != nil {
				zw.Close()
				out.Close()
				_ = os.Remove(tmp)
				return err
			}
			if _, err := w.Write(payload); err != nil {
				zw.Close()
				out.Close()
				_ = os.Remove(tmp)
				return err
			}
			replaced = true
			continue
		}
		rc, err := f.Open()
		if err != nil {
			zw.Close()
			out.Close()
			_ = os.Remove(tmp)
			return err
		}
		hdr := f.FileHeader
		w, err := zw.CreateHeader(&hdr)
		if err != nil {
			rc.Close()
			zw.Close()
			out.Close()
			_ = os.Remove(tmp)
			return err
		}
		_, err = io.Copy(w, rc)
		rc.Close()
		if err != nil {
			zw.Close()
			out.Close()
			_ = os.Remove(tmp)
			return err
		}
	}
	if !replaced {
		payload, err := os.ReadFile(filePath)
		if err != nil {
			zw.Close()
			out.Close()
			_ = os.Remove(tmp)
			return err
		}
		w, err := zw.Create(entry)
		if err != nil {
			zw.Close()
			out.Close()
			_ = os.Remove(tmp)
			return err
		}
		if _, err := w.Write(payload); err != nil {
			zw.Close()
			out.Close()
			_ = os.Remove(tmp)
			return err
		}
	}
	if err := zw.Close(); err != nil {
		out.Close()
		_ = os.Remove(tmp)
		return err
	}
	if err := out.Close(); err != nil {
		_ = os.Remove(tmp)
		return err
	}
	return os.Rename(tmp, zipPath)
}

func copyFile(src, dest string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
		return err
	}
	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, in)
	return err
}

func copyTree(src, dest string) error {
	return filepath.WalkDir(src, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		target := filepath.Join(dest, rel)
		if d.IsDir() {
			return os.MkdirAll(target, 0o755)
		}
		return copyFile(path, target)
	})
}

func copyDirContents(src, dest string) error {
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}
	for _, e := range entries {
		from := filepath.Join(src, e.Name())
		to := filepath.Join(dest, e.Name())
		if e.IsDir() {
			if err := copyTree(from, to); err != nil {
				return err
			}
		} else {
			if err := copyFile(from, to); err != nil {
				return err
			}
		}
	}
	return nil
}

func zipDir(src, dest string) error {
	f, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer f.Close()
	zw := zip.NewWriter(f)
	defer zw.Close()
	base := filepath.Base(src)
	return filepath.WalkDir(src, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		name := filepath.ToSlash(filepath.Join(base, rel))
		if d.IsDir() {
			if rel == "." {
				return nil
			}
			_, err := zw.Create(name + "/")
			return err
		}
		w, err := zw.Create(name)
		if err != nil {
			return err
		}
		in, err := os.Open(path)
		if err != nil {
			return err
		}
		_, copyErr := io.Copy(w, in)
		in.Close()
		return copyErr
	})
}

func tarGzDir(src, dest string) error {
	f, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer f.Close()
	gz := gzip.NewWriter(f)
	defer gz.Close()
	tw := tar.NewWriter(gz)
	defer tw.Close()
	base := filepath.Base(src)
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		hdr, err := tar.FileInfoHeader(info, "")
		if err != nil {
			return err
		}
		hdr.Name = filepath.ToSlash(filepath.Join(base, rel))
		if err := tw.WriteHeader(hdr); err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		in, err := os.Open(path)
		if err != nil {
			return err
		}
		_, copyErr := io.Copy(tw, in)
		in.Close()
		return copyErr
	})
}
