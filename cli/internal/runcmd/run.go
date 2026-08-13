package runcmd

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/MiaGobble/loveu/cli/internal/project"
	rt "github.com/MiaGobble/loveu/cli/internal/runtime"
)

type Options struct {
	Offline bool
	Args    []string
}

func Run(opts Options) error {
	proj, err := project.FindAndLoad(".")
	if err != nil {
		return err
	}
	preferConsole := proj.Console || runtime.GOOS == "windows"
	exe, err := rt.ResolveBinary(proj.EngineVersion, rt.HostPlatform(), preferConsole, opts.Offline)
	if err != nil {
		return err
	}
	got, out, err := rt.ProbeVersion(exe)
	if err != nil {
		return fmt.Errorf("probe runtime: %w\n%s", err, out)
	}
	if got != proj.EngineVersion {
		return fmt.Errorf("loveu.toml engine_version %q does not match running loveu %q", proj.EngineVersion, got)
	}

	args := append([]string{proj.Root}, opts.Args...)
	cmd := exec.Command(exe, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = proj.Root
	return cmd.Run()
}
