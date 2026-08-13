package versioncmd

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/MiaGobble/loveu/cli/internal/project"
	rt "github.com/MiaGobble/loveu/cli/internal/runtime"
	"github.com/MiaGobble/loveu/cli/internal/version"
)

func PrintCLI() {
	fmt.Printf("loveu %s\n", version.Version)
}

func Check(installMissing bool) error {
	PrintCLI()
	proj, err := project.FindAndLoad(".")
	if err != nil {
		if err == project.ErrNotFound {
			fmt.Println("project: (none)")
			cached, _ := rt.ListCached()
			if len(cached) == 0 {
				fmt.Println("runtimes: (none)")
			} else {
				fmt.Println("runtimes:")
				for _, c := range cached {
					fmt.Printf("  %s\n", c)
				}
			}
			return nil
		}
		return err
	}

	plat := rt.HostPlatform()
	fmt.Printf("project:  %s %s\n", proj.Name, proj.Version)
	fmt.Printf("pin:      %s\n", proj.EngineVersion)

	// On Windows prefer lovec.exe so --version stdout is captured (love.exe is GUI).
	preferConsole := runtime.GOOS == "windows"
	exe, err := resolveProbeBinary(proj.EngineVersion, plat, preferConsole)
	if err != nil {
		if installMissing {
			dir, instErr := rt.Install(rt.InstallOptions{Version: proj.EngineVersion, Platform: plat})
			if instErr != nil {
				fmt.Printf("runtime:  missing (%s)\n", plat)
				fmt.Printf("status:   error\n")
				return instErr
			}
			exe, err = rt.FindBinary(dir, plat, preferConsole)
			if err != nil {
				return err
			}
		} else {
			fmt.Printf("runtime:  missing (%s)\n", plat)
			fmt.Printf("status:   missing\n")
			fmt.Fprintf(os.Stderr, "hint: run `loveu version install` or `loveu version --install`\n")
			return fmt.Errorf("runtime %s/%s not installed", proj.EngineVersion, plat)
		}
	}

	got, _, err := rt.ProbeVersion(exe)
	if err != nil {
		fmt.Printf("runtime:  %s (%s)\n", proj.EngineVersion, plat)
		fmt.Printf("status:   error\n")
		return err
	}
	fmt.Printf("runtime:  %s  (%s)\n", got, plat)
	if got != proj.EngineVersion {
		fmt.Printf("status:   mismatch\n")
		return fmt.Errorf("engine_version %q does not match running loveu %q", proj.EngineVersion, got)
	}
	fmt.Printf("status:   ok\n")
	return nil
}

func resolveProbeBinary(ver string, plat rt.Platform, preferConsole bool) (string, error) {
	dir, err := rt.CachePath(ver, plat)
	if err != nil {
		return "", err
	}
	return rt.FindBinary(dir, plat, preferConsole)
}

func List() error {
	cached, err := rt.ListCached()
	if err != nil {
		return err
	}
	if len(cached) == 0 {
		fmt.Println("(no cached runtimes)")
		return nil
	}
	for _, c := range cached {
		fmt.Println(c)
	}
	return nil
}

type InstallFlags struct {
	Version  string
	Platform string
	From     string
}

func Install(flags InstallFlags) error {
	ver := flags.Version
	if ver == "" {
		if proj, err := project.FindAndLoad("."); err == nil {
			ver = proj.EngineVersion
		}
	}
	plat := rt.HostPlatform()
	if flags.Platform != "" {
		p, err := rt.ParsePlatform(flags.Platform)
		if err != nil {
			return err
		}
		plat = p
	}
	dir, err := rt.Install(rt.InstallOptions{
		Version:  strings.TrimPrefix(ver, "v"),
		Platform: plat,
		From:     flags.From,
	})
	if err != nil {
		return err
	}
	fmt.Printf("installed %s/%s -> %s\n", ver, plat, dir)
	return nil
}
