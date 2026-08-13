package runtime

import (
	"archive/zip"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"time"
)

const (
	DefaultReleasesRepo = "MiaGobble/loveu"
	EnvHome             = "LOVEU_HOME"
	EnvReleasesRepo     = "LOVEU_RELEASES_REPO"
)

var versionLine = regexp.MustCompile(`(?i)loveu\s+([0-9]+\.[0-9]+\.[0-9]+)`)

// Platform is an OS/arch key used in asset names and cache paths.
type Platform string

const (
	WindowsAMD64 Platform = "windows-x86_64"
	WindowsARM64 Platform = "windows-arm64"
	MacOS        Platform = "macos"
	LinuxAMD64   Platform = "linux-x86_64"
	LinuxARM64   Platform = "linux-aarch64"
	Android      Platform = "android"
	IOS          Platform = "ios"
)

func HostPlatform() Platform {
	switch runtime.GOOS {
	case "windows":
		if runtime.GOARCH == "arm64" {
			return WindowsARM64
		}
		return WindowsAMD64
	case "darwin":
		return MacOS
	case "linux":
		if runtime.GOARCH == "arm64" {
			return LinuxARM64
		}
		return LinuxAMD64
	default:
		return Platform(runtime.GOOS + "-" + runtime.GOARCH)
	}
}

func ParsePlatform(s string) (Platform, error) {
	switch Platform(s) {
	case WindowsAMD64, WindowsARM64, MacOS, LinuxAMD64, LinuxARM64, Android, IOS:
		return Platform(s), nil
	default:
		return "", fmt.Errorf("unknown platform %q (want windows-x86_64, windows-arm64, macos, linux-x86_64, linux-aarch64, android, ios)", s)
	}
}

func HomeDir() (string, error) {
	if v := os.Getenv(EnvHome); v != "" {
		return v, nil
	}
	switch runtime.GOOS {
	case "windows":
		base := os.Getenv("LOCALAPPDATA")
		if base == "" {
			home, err := os.UserHomeDir()
			if err != nil {
				return "", err
			}
			base = filepath.Join(home, "AppData", "Local")
		}
		return filepath.Join(base, "loveu"), nil
	default:
		if xdg := os.Getenv("XDG_DATA_HOME"); xdg != "" {
			return filepath.Join(xdg, "loveu"), nil
		}
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		return filepath.Join(home, ".local", "share", "loveu"), nil
	}
}

func RuntimesDir() (string, error) {
	home, err := HomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, "runtimes"), nil
}

func CachePath(ver string, plat Platform) (string, error) {
	base, err := RuntimesDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(base, ver, string(plat)), nil
}

func ReleasesRepo() string {
	if v := os.Getenv(EnvReleasesRepo); v != "" {
		return v
	}
	return DefaultReleasesRepo
}

func AssetName(ver string, plat Platform) string {
	switch plat {
	case Android:
		return fmt.Sprintf("loveu-engine-%s-android.apk", ver)
	case IOS:
		return fmt.Sprintf("loveu-engine-%s-ios.zip", ver)
	case MacOS:
		return fmt.Sprintf("loveu-engine-%s-macos.zip", ver)
	case LinuxAMD64, LinuxARM64:
		return fmt.Sprintf("loveu-engine-%s-%s.AppImage", ver, plat)
	default:
		return fmt.Sprintf("loveu-engine-%s-%s.zip", ver, plat)
	}
}

// Installed reports whether a usable runtime exists for ver/plat.
func Installed(ver string, plat Platform) (string, bool) {
	dir, err := CachePath(ver, plat)
	if err != nil {
		return "", false
	}
	exe, err := FindBinary(dir, plat, false)
	if err != nil {
		return "", false
	}
	return exe, true
}

// FindBinary locates love/lovec (or app binary) inside a cache dir.
func FindBinary(dir string, plat Platform, preferConsole bool) (string, error) {
	switch plat {
	case WindowsAMD64, WindowsARM64:
		names := []string{"love.exe", "lovec.exe"}
		if preferConsole {
			names = []string{"lovec.exe", "love.exe"}
		}
		return findNamed(dir, names)
	case MacOS:
		candidates := []string{
			filepath.Join(dir, "love.app", "Contents", "MacOS", "love"),
			filepath.Join(dir, "Contents", "MacOS", "love"),
		}
		for _, c := range candidates {
			if st, err := os.Stat(c); err == nil && !st.IsDir() {
				return c, nil
			}
		}
		return findNamed(dir, []string{"love"})
	case LinuxAMD64, LinuxARM64:
		candidates := []string{
			filepath.Join(dir, "bin", "love"),
			filepath.Join(dir, "squashfs-root", "bin", "love"),
			filepath.Join(dir, "love"),
		}
		for _, c := range candidates {
			if st, err := os.Stat(c); err == nil && !st.IsDir() {
				return c, nil
			}
		}
		// AppImage itself
		matches, _ := filepath.Glob(filepath.Join(dir, "*.AppImage"))
		if len(matches) > 0 {
			return matches[0], nil
		}
		return "", fmt.Errorf("no love binary in %s", dir)
	case Android:
		matches, _ := filepath.Glob(filepath.Join(dir, "*.apk"))
		if len(matches) > 0 {
			return matches[0], nil
		}
		apk := filepath.Join(dir, "loveu.apk")
		if st, err := os.Stat(apk); err == nil && !st.IsDir() {
			return apk, nil
		}
		return "", fmt.Errorf("no android apk in %s", dir)
	case IOS:
		return dir, nil
	default:
		return "", fmt.Errorf("unsupported platform %s", plat)
	}
}

func findNamed(dir string, names []string) (string, error) {
	var found string
	_ = filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}
		base := strings.ToLower(d.Name())
		for _, n := range names {
			if base == strings.ToLower(n) {
				found = path
				return filepath.SkipAll
			}
		}
		return nil
	})
	if found == "" {
		return "", fmt.Errorf("could not find %v under %s", names, dir)
	}
	return found, nil
}

// ProbeVersion runs love --version and parses the loveu semver.
func ProbeVersion(exe string) (string, string, error) {
	cmd := exec.Command(exe, "--version")
	out, err := cmd.CombinedOutput()
	text := string(out)
	if err != nil && text == "" {
		return "", text, err
	}
	m := versionLine.FindStringSubmatch(text)
	if m == nil {
		return "", text, fmt.Errorf("could not parse loveu version from: %s", strings.TrimSpace(text))
	}
	return m[1], text, nil
}

// ListCached returns installed runtime dirs as "ver/platform".
func ListCached() ([]string, error) {
	base, err := RuntimesDir()
	if err != nil {
		return nil, err
	}
	var out []string
	vers, err := os.ReadDir(base)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	for _, v := range vers {
		if !v.IsDir() {
			continue
		}
		plats, err := os.ReadDir(filepath.Join(base, v.Name()))
		if err != nil {
			continue
		}
		for _, p := range plats {
			if p.IsDir() {
				out = append(out, filepath.ToSlash(filepath.Join(v.Name(), p.Name())))
			}
		}
	}
	return out, nil
}

type releaseAsset struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
}

type releaseInfo struct {
	TagName string         `json:"tag_name"`
	Assets  []releaseAsset `json:"assets"`
}

func httpClient() *http.Client {
	return &http.Client{Timeout: 5 * time.Minute}
}

func fetchRelease(tag string) (*releaseInfo, error) {
	repo := ReleasesRepo()
	url := fmt.Sprintf("https://api.github.com/repos/%s/releases/tags/%s", repo, tag)
	if tag == "" || tag == "latest" {
		url = fmt.Sprintf("https://api.github.com/repos/%s/releases/latest", repo)
	}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", "loveu-cli")
	resp, err := httpClient().Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("no GitHub release %s (repo %s)", tag, repo)
	}
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		return nil, fmt.Errorf("GitHub API %s: %s", resp.Status, strings.TrimSpace(string(body)))
	}
	var info releaseInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return nil, err
	}
	return &info, nil
}

func findAsset(info *releaseInfo, name string) (*releaseAsset, error) {
	for i := range info.Assets {
		if info.Assets[i].Name == name {
			return &info.Assets[i], nil
		}
	}
	return nil, fmt.Errorf("no GitHub release %s (asset %s)", info.TagName, name)
}

func downloadFile(url, dest string) error {
	if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", "loveu-cli")
	resp, err := httpClient().Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download %s: %s", url, resp.Status)
	}
	tmp := dest + ".partial"
	f, err := os.Create(tmp)
	if err != nil {
		return err
	}
	_, copyErr := io.Copy(f, resp.Body)
	closeErr := f.Close()
	if copyErr != nil {
		_ = os.Remove(tmp)
		return copyErr
	}
	if closeErr != nil {
		_ = os.Remove(tmp)
		return closeErr
	}
	return os.Rename(tmp, dest)
}

func fileSHA256(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

func parseSHA256SUMS(data []byte) map[string]string {
	out := map[string]string{}
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) >= 2 {
			out[parts[len(parts)-1]] = strings.ToLower(parts[0])
		}
	}
	return out
}

func unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()
	for _, f := range r.File {
		target := filepath.Join(dest, f.Name)
		// zip-slip guard
		if !strings.HasPrefix(filepath.Clean(target)+string(os.PathSeparator), filepath.Clean(dest)+string(os.PathSeparator)) &&
			filepath.Clean(target) != filepath.Clean(dest) {
			return fmt.Errorf("illegal path in zip: %s", f.Name)
		}
		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(target, 0o755); err != nil {
				return err
			}
			continue
		}
		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
			return err
		}
		rc, err := f.Open()
		if err != nil {
			return err
		}
		out, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, f.Mode())
		if err != nil {
			rc.Close()
			return err
		}
		_, err = io.Copy(out, rc)
		out.Close()
		rc.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

func extractAppImage(src, dest string) error {
	if err := os.MkdirAll(dest, 0o755); err != nil {
		return err
	}
	// Prefer copying AppImage into dest; try --appimage-extract when executable.
	copied := filepath.Join(dest, filepath.Base(src))
	if err := copyFile(src, copied); err != nil {
		return err
	}
	_ = os.Chmod(copied, 0o755)
	cmd := exec.Command(copied, "--appimage-extract")
	cmd.Dir = dest
	if err := cmd.Run(); err != nil {
		// Keep AppImage as-is; fuse/extract may fail off-Linux.
		return nil
	}
	return nil
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

// InstallOptions controls version install.
type InstallOptions struct {
	Version  string
	Platform Platform
	From     string // local zip/dir/apk
	Offline  bool
}

// EnsureInstalled installs if missing (unless Offline).
func EnsureInstalled(opts InstallOptions) (string, error) {
	if opts.Platform == "" {
		opts.Platform = HostPlatform()
	}
	if opts.Version == "" {
		return "", errors.New("engine version is required")
	}
	if dir, ok := Installed(opts.Version, opts.Platform); ok {
		return dir, nil
	}
	if opts.Offline {
		return "", fmt.Errorf("runtime %s/%s not cached (offline)", opts.Version, opts.Platform)
	}
	return Install(opts)
}

// Install downloads or copies a runtime into the cache. Returns cache dir.
func Install(opts InstallOptions) (string, error) {
	if opts.Platform == "" {
		opts.Platform = HostPlatform()
	}
	if opts.Version == "" && opts.From == "" {
		info, err := fetchRelease("latest")
		if err != nil {
			return "", err
		}
		opts.Version = strings.TrimPrefix(info.TagName, "v")
	}
	if opts.Version == "" {
		return "", errors.New("engine version is required")
	}

	dest, err := CachePath(opts.Version, opts.Platform)
	if err != nil {
		return "", err
	}
	_ = os.RemoveAll(dest)
	if err := os.MkdirAll(dest, 0o755); err != nil {
		return "", err
	}

	var archive string
	cleanup := false
	if opts.From != "" {
		st, err := os.Stat(opts.From)
		if err != nil {
			return "", err
		}
		if st.IsDir() {
			if err := copyTree(opts.From, dest); err != nil {
				return "", err
			}
			return dest, nil
		}
		archive = opts.From
	} else {
		if opts.Offline {
			return "", errors.New("cannot download in offline mode")
		}
		tag := "v" + opts.Version
		info, err := fetchRelease(tag)
		if err != nil {
			return "", err
		}
		assetName := AssetName(opts.Version, opts.Platform)
		asset, err := findAsset(info, assetName)
		if err != nil {
			return "", err
		}
		tmpDir, err := os.MkdirTemp("", "loveu-dl-*")
		if err != nil {
			return "", err
		}
		defer os.RemoveAll(tmpDir)
		archive = filepath.Join(tmpDir, assetName)
		fmt.Fprintf(os.Stderr, "downloading %s...\n", assetName)
		if err := downloadFile(asset.BrowserDownloadURL, archive); err != nil {
			return "", err
		}
		// Optional checksum
		if sumsAsset, err := findAsset(info, "SHA256SUMS.txt"); err == nil {
			sumsPath := filepath.Join(tmpDir, "SHA256SUMS.txt")
			if err := downloadFile(sumsAsset.BrowserDownloadURL, sumsPath); err == nil {
				data, _ := os.ReadFile(sumsPath)
				sums := parseSHA256SUMS(data)
				if want, ok := sums[assetName]; ok {
					got, err := fileSHA256(archive)
					if err != nil {
						return "", err
					}
					if !strings.EqualFold(got, want) {
						return "", fmt.Errorf("SHA-256 mismatch for %s: got %s want %s", assetName, got, want)
					}
				}
			}
		}
		cleanup = true
		_ = cleanup
	}

	lower := strings.ToLower(archive)
	switch {
	case strings.HasSuffix(lower, ".apk"):
		if err := copyFile(archive, filepath.Join(dest, filepath.Base(archive))); err != nil {
			return "", err
		}
	case strings.HasSuffix(lower, ".appimage"):
		if err := extractAppImage(archive, dest); err != nil {
			return "", err
		}
	case strings.HasSuffix(lower, ".zip"):
		if err := unzip(archive, dest); err != nil {
			return "", err
		}
	default:
		// treat as opaque file
		if err := copyFile(archive, filepath.Join(dest, filepath.Base(archive))); err != nil {
			return "", err
		}
	}
	return dest, nil
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

// ResolveBinary returns the love binary path for a version/platform, installing if needed.
func ResolveBinary(ver string, plat Platform, preferConsole, offline bool) (string, error) {
	dir, err := EnsureInstalled(InstallOptions{Version: ver, Platform: plat, Offline: offline})
	if err != nil {
		return "", err
	}
	return FindBinary(dir, plat, preferConsole)
}
