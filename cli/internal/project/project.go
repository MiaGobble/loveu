package project

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"
)

const ManifestName = "loveu.toml"

var (
	ErrNotFound = errors.New("loveu.toml not found (walked up from current directory)")
)

// Project is the required loveu.toml surface the CLI needs.
type Project struct {
	Root           string
	Name           string `toml:"name"`
	Version        string `toml:"version"`
	EngineVersion  string `toml:"engine_version"`
	CodeRoot       string `toml:"code_root"`
	Console        bool   `toml:"console"`
	Title          string `toml:"title"`
	ManifestPath   string
}

type rawManifest struct {
	Name          string `toml:"name"`
	Version       string `toml:"version"`
	EngineVersion string `toml:"engine_version"`
	CodeRoot      string `toml:"code_root"`
	Console       bool   `toml:"console"`
	Title         string `toml:"title"`
}

// FindWalkUp searches cwd and parents for loveu.toml.
func FindWalkUp(start string) (string, error) {
	dir, err := filepath.Abs(start)
	if err != nil {
		return "", err
	}
	for {
		candidate := filepath.Join(dir, ManifestName)
		if st, err := os.Stat(candidate); err == nil && !st.IsDir() {
			return candidate, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", ErrNotFound
		}
		dir = parent
	}
}

// Load reads and validates a loveu.toml at path.
func Load(path string) (*Project, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var raw rawManifest
	if err := toml.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("invalid loveu.toml: %w", err)
	}
	for _, pair := range []struct {
		key, val string
	}{
		{"name", raw.Name},
		{"version", raw.Version},
		{"engine_version", raw.EngineVersion},
		{"code_root", raw.CodeRoot},
	} {
		if pair.val == "" {
			return nil, fmt.Errorf("loveu.toml is missing required field %q", pair.key)
		}
	}
	codeRoot := raw.CodeRoot
	if codeRoot == "" {
		codeRoot = "."
	}
	root := filepath.Dir(path)
	return &Project{
		Root:          root,
		Name:          raw.Name,
		Version:       raw.Version,
		EngineVersion: raw.EngineVersion,
		CodeRoot:      codeRoot,
		Console:       raw.Console,
		Title:         raw.Title,
		ManifestPath:  path,
	}, nil
}

// FindAndLoad walks up from start and loads the manifest.
func FindAndLoad(start string) (*Project, error) {
	path, err := FindWalkUp(start)
	if err != nil {
		return nil, err
	}
	return Load(path)
}

// MainLuau returns the expected main.luau path under code_root.
func (p *Project) MainLuau() string {
	return filepath.Join(p.Root, p.CodeRoot, "main.luau")
}
