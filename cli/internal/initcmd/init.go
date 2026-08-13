package initcmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/MiaGobble/loveu/cli/internal/runtime"
	"github.com/MiaGobble/loveu/cli/internal/version"
)

type Options struct {
	Dir   string
	Name  string
	Force bool
}

func Run(opts Options) error {
	dir := opts.Dir
	if dir == "" {
		dir = "."
	}
	abs, err := filepath.Abs(dir)
	if err != nil {
		return err
	}
	name := opts.Name
	if name == "" {
		name = filepath.Base(abs)
		if name == "." || name == string(filepath.Separator) || name == "" {
			name = "mygame"
		}
	}
	name = sanitizeName(name)

	if err := os.MkdirAll(abs, 0o755); err != nil {
		return err
	}
	entries, err := os.ReadDir(abs)
	if err != nil {
		return err
	}
	if len(entries) > 0 && !opts.Force {
		return fmt.Errorf("directory %s is not empty (use --force to overwrite scaffold files)", abs)
	}

	engine := version.Version
	if cached, err := runtime.ListCached(); err == nil && len(cached) > 0 {
		// Prefer latest-looking first entry's version segment
		parts := strings.Split(cached[len(cached)-1], "/")
		if len(parts) > 0 && parts[0] != "" {
			engine = parts[0]
		}
	}

	toml := fmt.Sprintf(`name = %q
version = "0.0.1"
engine_version = %q
code_root = "src"

[window]
width = 800
height = 600
`, name, engine)

	mainLuau := `function love.draw()
	love.graphics.print("Hello from Luau!", 100, 100)
end
`

	gitignore := `dist/
.DS_Store
*.love
`

	srcDir := filepath.Join(abs, "src")
	if err := os.MkdirAll(srcDir, 0o755); err != nil {
		return err
	}
	if err := writeFile(filepath.Join(abs, "loveu.toml"), toml); err != nil {
		return err
	}
	if err := writeFile(filepath.Join(srcDir, "main.luau"), mainLuau); err != nil {
		return err
	}
	if err := writeFile(filepath.Join(abs, ".gitignore"), gitignore); err != nil {
		return err
	}
	fmt.Printf("Created loveu project %q in %s\n", name, abs)
	fmt.Println("Next: loveu run")
	return nil
}

func writeFile(path, contents string) error {
	return os.WriteFile(path, []byte(contents), 0o644)
}

func sanitizeName(s string) string {
	s = strings.TrimSpace(s)
	s = strings.ReplaceAll(s, " ", "-")
	return s
}
