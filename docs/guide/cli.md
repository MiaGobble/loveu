# CLI

loveu development is centered on the **`loveu`** command, not a raw `love` binary.

Install a release installer from [GitHub Releases](https://github.com/MiaGobble/loveu/releases) (Windows setup.exe, macOS `.pkg`, Linux `.deb`), or build from this repo:

```
cd cli
go build -o loveu ./cmd/loveu
```

## Commands

| Command | Purpose |
|---------|---------|
| `loveu init [dir]` | Scaffold `loveu.toml` + `src/main.luau` |
| `loveu run` | Download the pinned engine if needed, then run the project |
| `loveu build <target>` | Package for `love`, `windows`, `macos`, `linux`, `android`, or `ios` |
| `loveu version` | Check `engine_version` against the cached runtime |
| `loveu version install [ver]` | Download an engine pack into the local cache |
| `loveu version list` | List cached engines |

`engine_version` in `loveu.toml` is the selected engine. There is no separate global pin.

## Engine download

Packs come from GitHub Releases (`v<engine_version>`). Override with `LOVEU_RELEASES_REPO`. Cache root defaults to `%LOCALAPPDATA%\loveu` (Windows) or `~/.local/share/loveu` (override with `LOVEU_HOME`).

Use `--from <zip-or-dir>` for a local CI artifact, or `--offline` on `run` / `build` to forbid network.

## Build targets

Desktop fuse/package, Android APK inject, and iOS (macOS + Xcode only). Consoles are not supported — upstream LÖVE never shipped Switch/Xbox/PlayStation backends.
