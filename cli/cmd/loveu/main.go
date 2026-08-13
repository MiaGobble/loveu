package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/MiaGobble/loveu/cli/internal/buildcmd"
	"github.com/MiaGobble/loveu/cli/internal/initcmd"
	"github.com/MiaGobble/loveu/cli/internal/runcmd"
	"github.com/MiaGobble/loveu/cli/internal/version"
	"github.com/MiaGobble/loveu/cli/internal/versioncmd"
)

func main() {
	os.Exit(run(os.Args[1:]))
}

func run(args []string) int {
	if len(args) == 0 {
		printHelp()
		return 0
	}
	switch args[0] {
	case "-h", "--help", "help":
		printHelp()
		return 0
	case "-v", "--version":
		versioncmd.PrintCLI()
		return 0
	case "version":
		return cmdVersion(args[1:])
	case "init":
		return cmdInit(args[1:])
	case "run":
		return cmdRun(args[1:])
	case "build":
		return cmdBuild(args[1:])
	default:
		fmt.Fprintf(os.Stderr, "unknown command %q\n\n", args[0])
		printHelp()
		return 2
	}
}

func printHelp() {
	fmt.Printf(`loveu %s — develop and package loveu games

Usage:
  loveu <command> [options]

Commands:
  version                 Check project pin vs cached engine
  version install [ver]   Download engine into local cache
  version list            List cached engines
  init [dir]              Create a blank project
  run [--] [args...]      Run the project (downloads engine if needed)
  build <target>          Package: love|windows|macos|linux|android|ios

Global:
  -h, --help              Show help
  -v, --version           Print CLI version

Environment:
  LOVEU_HOME              Cache root (default: platform data dir)
  LOVEU_RELEASES_REPO     GitHub repo for engine packs (default: MiaGobble/loveu)
`, version.Version)
}

func cmdVersion(args []string) int {
	installFlag := false
	rest := []string{}
	for i := 0; i < len(args); i++ {
		a := args[i]
		switch a {
		case "install":
			return cmdVersionInstall(args[i+1:])
		case "list":
			if err := versioncmd.List(); err != nil {
				fmt.Fprintln(os.Stderr, err)
				return 1
			}
			return 0
		case "--install":
			installFlag = true
		case "-h", "--help":
			fmt.Println(`Usage:
  loveu version [--install]
  loveu version install [ver] [--platform <plat>] [--from <path>]
  loveu version list`)
			return 0
		default:
			rest = append(rest, a)
		}
	}
	_ = rest
	if err := versioncmd.Check(installFlag); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	return 0
}

func cmdVersionInstall(args []string) int {
	flags := versioncmd.InstallFlags{}
	for i := 0; i < len(args); i++ {
		a := args[i]
		switch a {
		case "--platform":
			if i+1 >= len(args) {
				fmt.Fprintln(os.Stderr, "--platform requires a value")
				return 2
			}
			i++
			flags.Platform = args[i]
		case "--from":
			if i+1 >= len(args) {
				fmt.Fprintln(os.Stderr, "--from requires a value")
				return 2
			}
			i++
			flags.From = args[i]
		case "-h", "--help":
			fmt.Println("Usage: loveu version install [ver] [--platform <plat>] [--from <path>]")
			return 0
		default:
			if flags.Version == "" && !hasPrefix(a, "-") {
				flags.Version = a
			} else {
				fmt.Fprintf(os.Stderr, "unknown option %q\n", a)
				return 2
			}
		}
	}
	if err := versioncmd.Install(flags); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	return 0
}

func cmdInit(args []string) int {
	opts := initcmd.Options{}
	for i := 0; i < len(args); i++ {
		a := args[i]
		switch a {
		case "--force":
			opts.Force = true
		case "--name":
			if i+1 >= len(args) {
				fmt.Fprintln(os.Stderr, "--name requires a value")
				return 2
			}
			i++
			opts.Name = args[i]
		case "-h", "--help":
			fmt.Println("Usage: loveu init [dir] [--name <name>] [--force]")
			return 0
		default:
			if opts.Dir == "" && !hasPrefix(a, "-") {
				opts.Dir = a
			} else {
				fmt.Fprintf(os.Stderr, "unknown option %q\n", a)
				return 2
			}
		}
	}
	if err := initcmd.Run(opts); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	return 0
}

func cmdRun(args []string) int {
	opts := runcmd.Options{}
	gameArgs := []string{}
	passthrough := false
	for i := 0; i < len(args); i++ {
		a := args[i]
		if passthrough {
			gameArgs = append(gameArgs, a)
			continue
		}
		switch a {
		case "--":
			passthrough = true
		case "--offline":
			opts.Offline = true
		case "-h", "--help":
			fmt.Println("Usage: loveu run [--offline] [--] [game args...]")
			return 0
		default:
			gameArgs = append(gameArgs, a)
		}
	}
	opts.Args = gameArgs
	if err := runcmd.Run(opts); err != nil {
		var exitErr *exitError
		if errors.As(err, &exitErr) {
			return exitErr.Code
		}
		// os/exec.ExitError
		type coder interface{ ExitCode() int }
		if c, ok := err.(coder); ok {
			return c.ExitCode()
		}
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	return 0
}

type exitError struct{ Code int }

func (e *exitError) Error() string { return fmt.Sprintf("exit %d", e.Code) }

func cmdBuild(args []string) int {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "Usage: loveu build <love|windows|macos|linux|android|ios> [--offline] [--keystore path]")
		return 2
	}
	opts := buildcmd.Options{Target: args[0]}
	for i := 1; i < len(args); i++ {
		a := args[i]
		switch a {
		case "--offline":
			opts.Offline = true
		case "--keystore":
			if i+1 >= len(args) {
				fmt.Fprintln(os.Stderr, "--keystore requires a value")
				return 2
			}
			i++
			opts.Keystore = args[i]
		case "-h", "--help":
			fmt.Println("Usage: loveu build <love|windows|macos|linux|android|ios> [--offline] [--keystore path]")
			return 0
		default:
			fmt.Fprintf(os.Stderr, "unknown option %q\n", a)
			return 2
		}
	}
	if err := buildcmd.Run(opts); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	return 0
}

func hasPrefix(s, p string) bool {
	return len(s) >= len(p) && s[:len(p)] == p
}
